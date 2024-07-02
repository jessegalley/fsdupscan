package sizetree_test

import (
	// "fmt"
	"testing"

	"github.com/jessegalley/fsdupscan/internal/sizetree"
	"github.com/stretchr/testify/assert"
)

// Test initialization of the SizeTree
func TestNewSizeTree(t *testing.T) {
  st := sizetree.New()
  
  assert.IsType(t, &sizetree.SizeTree{}, st)
}

// Test insertion of a SizeTreeEntry into the SizeTree 
// for both existing and non existing sizes.
func TestReplaceOrInsert(t *testing.T) {
  st := sizetree.New()

  e := sizetree.NewSizeTreeEntry(1000, []string{"/test/file/one"})
  node := st.ReplaceOrInsert(e) 
  assert.Nil(t, node)

  e2 := sizetree.NewSizeTreeEntry(1000, []string{"/test/file/two"})
  node2 := st.ReplaceOrInsert(e2)
  assert.Same(t, e, node2)
}

// Test getting of existing and non existing SizeTreeEntries from
// the SizeTree 
func TestGet(t *testing.T) {
  st := sizetree.New()

  e1 := sizetree.NewSizeTreeEntry(1000, []string{"/test/file/one"})
  item1 := st.ReplaceOrInsert(e1)
  assert.Nil(t, item1)

  e2 := sizetree.NewSizeTreeEntry(1000, nil)
  item2 := st.Get(e2)
  assert.Same(t, item2, e1)

  e3 := sizetree.NewSizeTreeEntry(1001, nil)
  item3 := st.Get(e3)
  assert.Nil(t, item3)
}

// Test getting of existing and non existing SizeTreeEntries by Size alone
func TestGetBySize(t *testing.T) {
  st := sizetree.New()

  e1 := sizetree.NewSizeTreeEntry(1000, []string{"/test/file/one"})
  item1 := st.ReplaceOrInsert(e1)
  assert.Nil(t, item1)
  
  item2 := st.GetBySize(1000)
  assert.Same(t, e1, item2)

  item3 := st.GetBySize(1001)
  assert.Nil(t, item3)
}

// Test appending an additional filename to a SizeTreeEntry
func TestAppend(t *testing.T) {
  // st := sizetree.New()

  e1 := sizetree.NewSizeTreeEntry(1000, []string{"/file/one"})
  // item1 := st.ReplaceOrInsert(e1)
  // assert.Nil(t, item1)

  e1.Append("/test/two")

  files := e1.Files()
  assert.Equal(t, len(files), 2)

  assert.Equal(t, files[0], "/file/one")
  assert.Equal(t, files[1], "/test/two")
}

// tests merging two SizeTreeEntry's files into one 
func TestMerge(t *testing.T) {
  e1 := sizetree.NewSizeTreeEntry(1000, []string{"/file/one"})
  e2 := sizetree.NewSizeTreeEntry(1000, []string{"/test/two"})

  e1.Merge(e2)
  files := e1.Files()
  assert.Equal(t, len(files), 2)

  assert.Equal(t, files[0], "/file/one")
  assert.Equal(t, files[1], "/test/two")
}

func TestMergOrInsert(t *testing.T) {
  st := sizetree.New()

  e1 := sizetree.NewSizeTreeEntry(1000, []string{"/file/one"})
  item1 := st.MergeOrInsert(e1)
  assert.Nil(t, item1)

  e2 := sizetree.NewSizeTreeEntry(1000, []string{"/test/two"})
  item2 := st.MergeOrInsert(e2)
  assert.IsType(t, &sizetree.SizeTreeEntry{}, item2)

  files := item2.Files()
  assert.Equal(t, len(files), 2)
  assert.Equal(t, files[0], "/file/one")
  assert.Equal(t, files[1], "/test/two")

  // fmt.Println(item2)
  // fmt.Println(st.GetBySize(1000))
}






