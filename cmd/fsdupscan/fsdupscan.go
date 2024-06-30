package main

import (
	"fmt"
	"log"
	"log/slog"
	"os"
	"sync"

	// "path/filepath"
	// "sync"

	// "time"

	// "github.com/spf13/pflag"
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
  setupCliArgs()
  setupLogger()
}

func main() { 
  slog.Info("main", "working", "working")
  slog.Debug("pos args", "args", positionalArgs)
  slog.Debug("main", "threads", viper.GetInt("threads"))

  // setup a waitgroup and a semaphore for limiting 
  // concurrent scan actions
  // var wg sync.WaitGroup
  // sem := make(chan struct{}, viper.GetInt("threads"))
  dirCh := make(chan string, 1)
  // fileCh := make(chan os.FileInfo, 1)
  
  // add all of the paths provided as positionalArgs to the directory channel
  // if there were no positional args provided, then assume CWD and add that
  // to the directory channel
  go func ()  {
    if len(positionalArgs) == 0 {
      dirCh <- "."
    } else {
      for _, path := range positionalArgs {
        dirCh <- path
      }
    }  
  }()
  
  var wg sync.WaitGroup
  wg.Add(1)
  go func ()  {
    defer wg.Done()
    for {
      select {
      case dir, ok := <- dirCh:
        //do something with dir 
        slog.Debug("dirCh processing", "dir", dir)
        if !ok {
          return
        }
      }
    }
  }()

  wg.Wait()
  // for each directory given in positional args 
  // validate that the path exists
  // then start a walk of that filepath
  // if there was no positional args specified, start walk on current directory 

  // for _, path := range positionalArgs {
  //   wg.Add(1)
  //   walkDir(path, &wg, sem)
  //   walkPath(path, &wg, )
  // }
  //
  // wg.Wait()
}

// func walkDir(dir string, wg *sync.WaitGroup, sem chan struct{}) {
//   defer wg.Done()
//
//
//   entries, err := os.ReadDir(dir)
//   if err != nil {
//     log.Println(err)
//     return
//   }
//
//   for _, entry := range entries {
//     path := filepath.Join(dir, entry.Name())
//     if entry.IsDir() {
//       // block if sem channel is full
//       slog.Debug("walkDir", "dir", path)
//       sem <- struct{}{}
//       slog.Debug("walkDir", "sem", path)
//       wg.Add(1)
//       go func(p string) {
//         defer func() { <-sem }()
//         walkDir(p, wg, sem)
//         slog.Debug("walkDir", "release", p)
//         // release the semaphore
//         // <-sem
//       }(path)
//     } else {
//       // process the file
//       // fmt.Println(path)
//     }
//   }
// }
  
