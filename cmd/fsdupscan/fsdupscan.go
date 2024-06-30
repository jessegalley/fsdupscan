package main

import (
	"fmt"
	"log"
	"log/slog"
	"os"
	// "time"

	// "github.com/spf13/pflag"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
  semVer = "0.1.0"
  progName = "fsdupscan"
)

//setupCliArgs wraps the various commandline arguments and options parsing
//and set up tasks for this program. It will also initiate the argparser 
//and handle basic housekeeping tasks like counting positional arguments 
//and handling arguments such as verson or help
func setupCliArgs () {

  viper.SetDefault("verbose", false)
  viper.SetDefault("debug", false)
  viper.SetDefault("tickrate", 1000)

  // viper.SetConfigName(progName)
  // viper.SetConfigType("yaml")
  // viper.AddConfigPath("/etc/"+progName+"/")
  // viper.AddConfigPath("/etc/")
  // viper.AddConfigPath("$HOME/."+progName+"/")
  // viper.AddConfigPath(".")
  // err := viper.ReadInConfig()
  // if err !=  nil {
  //   panic(fmt.Errorf("fatal error reading config file: %w", err))
  // }
  //
  // set up all commandline flags
  flag.BoolP("debug", "D", false, "print debug messages")
  flag.BoolP("verbose", "v", false, "be verbose")
  flag.BoolP("version", "V", false, "print version and exit")
  // flag.IntP("tickrate", "T", 1000, "service tickrate in ms")
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

  // make sure that an incorrect number of args wasn't provided
  expectedArgs := 0
  if len(flag.Args()) != expectedArgs {
    flag.Usage()
    os.Exit(2)
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
  // setup the ticker for the daemon
  // delay := time.Duration(viper.GetInt("tickrate") * int(time.Millisecond))
  // ticker := time.NewTicker(delay)
  // defer ticker.Stop()

  // main daemon loop
  // for {
  //   select {
  //   case <-ticker.C:
  //     //do work
  //     slog.Info("hello world")
  //     slog.Debug("hello Debug")
  //   }
  // }
}
