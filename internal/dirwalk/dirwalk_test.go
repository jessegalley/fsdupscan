package dirwalk_test

import (
	// "fmt"
	"os"
	"path/filepath"

	// "path/filepath"
	"sync"
	"testing"

	// "github.com/davecgh/go-spew/spew"
	"github.com/jessegalley/fsdupscan/internal/dirwalk"
	"github.com/stretchr/testify/assert"
)

func TestDirwalk(t *testing.T) {
  basePath := "../../testdata/"
  // basePath := "./testdata/"

  var wg sync.WaitGroup
  fileCh := make(chan os.DirEntry, 1)
  
  wg.Add(1)
  go func ()  {
    defer wg.Done()
    dirwalk.Walk(basePath, fileCh)
  }()

  files := []string{}
  go func ()  {
    for {
      select {
      case entry, ok := <-fileCh:
        if ok {
          files = append(files, filepath.Join(basePath,entry.Name())) 
          // fmt.Println(os.ModeSymlink)
          // spew.Dump(entry)
          // spew.Dump(entry.Type())
          // fmt.Println(entry.Type())
          // spew.Dump(entry.Info())
        } else {
          return
        }
      }
    }
  }()

  wg.Wait()
  assert.Equal(t, 18, len(files))

  // for _, entry := range files {
  //   // spew.Dump(entry)
  // }
}

func isSymlink(file os.DirEntry) bool {
  mode := file.Type()
  return mode & os.ModeSymlink != 0
}
