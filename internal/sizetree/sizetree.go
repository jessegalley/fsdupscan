package sizetree

import "github.com/google/btree"

// type SizeTree is a wrapper to google's btree, providing a
// simple interface to manage a btree of filepaths and sizes 
// for the purpose of the the fsdupscan utility only.
type SizeTree struct {
  btree *btree.BTree
}

// New() returns a pointer to a SizeTree struct, with an initialized btree
func New() *SizeTree  {
  return &SizeTree{
    btree: btree.New(2),
  }
}

// FileSizeEntry is a representation of a single file on the filesystem
// consisting of its name, path, and size.  These are the parameters needed
// to both identify a file and sort/search it by filesize.
type FileSizeEntry struct {
  Path  string 
  Name  string 
  Size  int64
}

func NewFileSizeEntry(path string, name string, size int64) *FileSizeEntry {
  return &FileSizeEntry{
    Path: path,
    Name: name,
    Size: size,
  }
}
