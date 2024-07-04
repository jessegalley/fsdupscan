package dirwalk_test

import (
	"testing"

	// "github.com/davecgh/go-spew/spew"
	"github.com/davecgh/go-spew/spew"
	"github.com/jessegalley/fsdupscan/internal/dirwalk"
	"github.com/stretchr/testify/assert"
)

func TestDirwalk(t *testing.T) {
  basePath := "../../testdata/"
  // basePath := "./testdata/"

  // var wg sync.WaitGroup
  // fileCh := make(chan os.DirEntry, 1)
  
  // wg.Add(1)
  // go func ()  {
  //   // defer wg.Done()
  //   followSym := false
  //   dirwalk.Walk(basePath, fileCh, followSym)
  // }()

  fileCh, wgWalk := dirwalk.Walk(basePath)
  files := []string{}
  go func ()  {
    for {
      select {
      case entry, ok := <-fileCh:
        if ok {
          // files = append(files, filepath.Join(basePath,entry.Name())) 
          files = append(files, entry) 
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

  wgWalk.Wait()
  assert.Equal(t, 18, len(files))

  spew.Dump(files)
  // for _, entry := range files {
  //   // spew.Dump(entry)
  // }
}


