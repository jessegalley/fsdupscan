package dirwalk

import (
	"log/slog"
	"os"
	"path/filepath"
)

func Walk(dir string, fileCh chan<- os.DirEntry) {
  entries, err := readDirRegular(dir)
  if err != nil {
    // we shouldn't get here
    // panic("something broke when reading dir in walk()")
    panic(err)
  }

  for _, entry := range entries {
    if entry.IsDir() {
      Walk(filepath.Join(dir, entry.Name()), fileCh)
    } else if entry.Type().IsRegular() {
      fileCh <- entry
    } else {
      slog.Debug("dirwalk::walk() unknown entry.Type()", "Type()", entry.Type())
      continue
    }
  }
}

func readDirRegular(dir string) ([]os.DirEntry, error) { 
  entries, err := os.ReadDir(dir)
  if err != nil {
    return nil, err
    // log.Fatal(err)
  }

  return entries, nil
}
