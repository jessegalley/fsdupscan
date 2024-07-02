package sizetree_test

import (
	"testing"

	"github.com/jessegalley/fsdupscan/internal/sizetree"
	"github.com/stretchr/testify/assert"
)


func TestNewSizeTree(t *testing.T) {
  st := sizetree.New()
  
  assert.IsType(t, &sizetree.SizeTree{}, st)
}
