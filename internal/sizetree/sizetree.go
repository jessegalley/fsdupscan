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

// SizeTreeEntry is a single node in the SizeTree btree structure. It is indexed
// on the size of the file, and contains a slice of files which share this size.
type SizeTreeEntry struct {
  Size  int64
  files   []string
}

// NewSizeTreeEntry() returns a pointer to a SizeTreeEntry, with initialized parameters.
func NewSizeTreeEntry (size int64, files []string) *SizeTreeEntry {
  return &SizeTreeEntry{
    Size: size,
    files: files,
  }
}

// Append() method appends a file to the SizeTreeEntry.
func (s *SizeTreeEntry) Append(file string) {
  s.files = append(s.files, file)
}

// Less() compares a SizeTreeEntry to another SizeTreeEntry 
// this is required to satisfy the btree.Item interface 
func (a SizeTreeEntry) Less(b btree.Item) bool {
	return a.Size < b.(*SizeTreeEntry).Size
}


