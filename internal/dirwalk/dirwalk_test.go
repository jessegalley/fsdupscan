package dirwalk_test

import (
	"os"
	"path/filepath"

	// "path/filepath"
	"sync"
	"testing"

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
        } else {
          return
        }
      }
    }
  }()

  wg.Wait()
  assert.Equal(t, len(files), 2)
}


