package main

import (
	"errors"
	"fmt"
	"log"
	"log/slog"
	"os"
	"runtime"
	"time"

	"github.com/jessegalley/fsdupscan/internal/dirwalk"
	"github.com/jessegalley/fsdupscan/internal/filechecksum"
	"github.com/jessegalley/fsdupscan/internal/sizetree"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
  semVer = "0.1.0"
  progName = "fsdupscan"
)

var ( 
  positionalArgs []string
)

//setupCliArgs wraps the various commandline arguments and options parsing
//and set up tasks for this program. It will also initiate the argparser 
//and handle basic housekeeping tasks like counting positional arguments 
//and handling arguments such as verson or help
func setupCliArgs () {
  viper.SetDefault("verbose", false)
  viper.SetDefault("debug", false)
  viper.SetDefault("tickrate", 1000)

  // set up all commandline flags
  flag.BoolP("debug", "D", false, "print debug messages")
  flag.BoolP("verbose", "v", false, "be verbose")
  flag.BoolP("version", "V", false, "print version and exit")
  // flag.IntP("tickrate", "T", 1000, "service tickrate in ms")
  flag.IntP("threads", "t", 10, "number of concurrent filesystem scan workers")
  flag.IntP("min", "m", 1000, "minimum file size to consider in bytes")
  flag.IntP("max", "M", 0, "maximum file size to consider in bytes (0 for no maximum)")
  // TODO: handle human readable size suffixes like "1K", "1M" etc...

  flag.CommandLine.MarkHidden("debug")
  flag.Usage =  func() {
    fmt.Fprintf(os.Stderr, "usage: %s [OPTS] \n", os.Args[0])
    fmt.Fprintf(os.Stderr, "\n")
    flag.PrintDefaults()
  }
  flag.Parse()

  viper.BindPFlags(flag.CommandLine)

  // if -v/--version is given, print version info and exit
  vflag, _  := flag.CommandLine.GetBool("version") 
  if vflag { 
    fmt.Println("v", semVer)
    os.Exit(1)
  }

  // we will accept any number of positional arguments here, as the 
  // user may specify any number of directories to include in the scan
  for _, arg := range flag.Args() {
    positionalArgs = append(positionalArgs, arg)
  }
}

// setupLogger wraps the various logger setup tasks for this program
func setupLogger () {
  if viper.GetBool("debug") {
    slog.SetLogLoggerLevel(slog.LevelDebug)
  }
  log.SetFlags(log.Ldate | log.Ltime | log.Lmsgprefix)
  log.SetPrefix(progName+": ")
}

func init() {
  numCPUs := runtime.NumCPU()
  runtime.GOMAXPROCS(numCPUs)
  setupCliArgs()
  setupLogger()
}

func main() { 
  slog.Debug("pos args", "args", positionalArgs)
  slog.Debug("main", "threads", viper.GetInt("threads"))

  minSize := viper.GetInt("min")

  // make sure all input arguments are acessible dirs 
  // TODO: come up with a way to handle files as args too 
  //       perhaps adding them directly to SizeTree
  _, err := validateStartingDirs(positionalArgs)
  if err != nil {
    slog.Error("error checking input dirs", "err", err)
    os.Exit(3)
  }

  slog.Info("starting scan", "dirs", positionalArgs)
  fileCh, wg := dirwalk.Walk(positionalArgs...)

  st := sizetree.New()


  var filesScanned int64
  var filesSkipped int64
  var filesCompared int64

  go func ()  {
    for {
      select {
      case file, ok := <- fileCh:
        if !ok {
          slog.Debug("fileCh closed")
          return
        }

        filesScanned++
        if file.Size < int64(minSize) {
          filesSkipped++
          slog.Debug("skipping too small", "file", file.Path, "size", file.Size)
          continue
        }
        filesCompared++
        f1 := sizetree.SizeTreeFile{Path: file.Path, Inode: file.Inode}
        e1 := sizetree.NewSizeTreeEntry(file.Size, []sizetree.SizeTreeFile{f1})
        item := st.MergeOrInsert(e1)
        if item != nil {
          // trigger a hash check!
          ste := st.GetBySize(file.Size)
          if ste == nil {
            // defensive programming error so if some race conidtion happens we dont 
            // panic with nil pointer dreferences 
            panic("nil SizeTreeEntry pointer after size collision, this shouldn't happen.")
          }
          stfs := ste.Files()
          if stfs == nil {
            // more defensive programming to catch race cond
            panic("nil/empty slice of sizetreefile after size collison, this shoudln't happen")
          }
          for _, stf := range stfs {
            checksum, err := filechecksum.CalculateChecksum(stf.Path)
            if err != nil {
              // not sure why hash would fail, maybe nil path?
              slog.Error("cant hash", "error", err)
              panic("hashing of file failed. file:"+file.Path)
            }
            // stf.SetHash(checksum)
            ste.AppendChecksum(checksum, &stf)
          }
          // slog.Info("main tree insert collision", "path", file.Path, "size", file.Size)
        } 
        if viper.GetBool("verbose") {
          fmt.Println(file) // what should --verbose actually print?
        }
      }
    }
  }()

  wg.Wait()


  time.Sleep(10 * time.Second)
  // slog.Debug("len(comparisons)", "len", len(comparisons))
  slog.Info("summary", "scanned", filesScanned, "skipped", filesSkipped, "compared", filesCompared)
  
}


func validateStartingDirs(dirs []string) (bool, error) {
  for _, dir := range dirs {
    info, err := os.Stat(dir)
    if err != nil {
      return false, errors.New("can't open path: "+ dir)
    }
    if !info.IsDir() {
      return false, errors.New("path isn't dir:"+ dir)
    }
  }

  return true, nil
}

