package dirwalk_test

import (
	"testing"

	"github.com/jessegalley/fsdupscan/internal/dirwalk"
	"github.com/stretchr/testify/assert"
)

func TestDirwalk(t *testing.T) {
  basePath := "../../testdata/"

  fileCh, wgWalk := dirwalk.Walk(basePath)
  files := []*dirwalk.WalkedFile{}
  go func ()  {
    for {
      select {
      case entry, ok := <-fileCh:
        if ok {
          files = append(files, entry) 
        } else {
          return
        }
      }
    }
  }()

  wgWalk.Wait()
  assert.Equal(t, 18, len(files))

}


