package sizetree

import "github.com/google/btree"


// type SizeTreeFile is a representation of a single file sorted 
// into the SizeTreeEntrys
type SizeTreeFile struct {
  Path string
  Inode int64
}

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

// ReplaceOrInsert inserts a SizeTreeEntry into the SizeTree, returning nil if 
// a SizeTreeEntry of that size did not exist, and returning the existing SizeTreeEntry
// if one of that size already existed.
func (s *SizeTree) ReplaceOrInsert (e *SizeTreeEntry) *SizeTreeEntry {
  node := s.btree.ReplaceOrInsert(e)
  if node == nil{
    return nil
  }

  return node.(*SizeTreeEntry)
}

// MergeOrInsert inserts he given SizeTreeEntry into the SizeTree 
// if an entry of that Size already exists, the two entries are merged
// returns nil if no merge occured
func (s *SizeTree) MergeOrInsert (e *SizeTreeEntry) *SizeTreeEntry {
  item := s.btree.ReplaceOrInsert(e)
  if item == nil {
    return nil
  }

  // item and e are different, meaning the Size already existed in the tree 
  // merge the two SizeTreeEntries
  item.(*SizeTreeEntry).Merge(e)

  // re-insert the merged entry 
  item2 := s.btree.ReplaceOrInsert(item)
  if item2 ==  nil {
    // we should never get here 
    panic("something when horribly wrong trying to merge an entry")
  }

  // finally return the merged entry back to the caller
  return item.(*SizeTreeEntry)
}

// Get returns a pointer to a SizeTreeEntry of matching Size from the SizeTree, 
// if no entry of this size  exists then it returns nil.
func (s *SizeTree) Get(e *SizeTreeEntry) *SizeTreeEntry {
  item := s.btree.Get(e)
  if item == nil {
    return nil
  }

  return item.(*SizeTreeEntry)
}

// GetBySize returns a pointer to a SizeTreeEntry of matching Size from the SizeTree,
// if no entry of this size  exists then it returns nil.
func (s *SizeTree) GetBySize(size int64) *SizeTreeEntry {
  e := NewSizeTreeEntry(size, nil)
  item := s.btree.Get(e)
  if item == nil {
    return nil
  }

  return item.(*SizeTreeEntry)
}

// SizeTreeEntry is a single node in the SizeTree btree structure. It is indexed
// on the size of the file, and contains a slice of files which share this size.
type SizeTreeEntry struct {
  Size  int64
  files   []SizeTreeFile
}

// NewSizeTreeEntry() returns a pointer to a SizeTreeEntry, with initialized parameters.
func NewSizeTreeEntry (size int64, files []SizeTreeFile) *SizeTreeEntry {
  return &SizeTreeEntry{
    Size: size,
    files: files,
  }
}

// Append() method appends a file to the SizeTreeEntry.
func (s *SizeTreeEntry) Append(file SizeTreeFile) {
  s.files = append(s.files, file)
}

// Merge method appends all files from the other SizeTreeEntry into this one.
func (s *SizeTreeEntry) Merge(other *SizeTreeEntry) {
  otherFiles := other.Files() 
  if otherFiles == nil {
    return
  }

  s.files = append(s.files, otherFiles...)
}

// Files retuns a []string of file paths associated with this SizeTreeEntry 
func (s *SizeTreeEntry) Files() []SizeTreeFile{
  if len(s.files) == 0 {
    return nil
  } else if s.files == nil {
    return nil
  }

  return s.files
}

// Less() compares a SizeTreeEntry to another SizeTreeEntry 
// this is required to satisfy the btree.Item interface 
func (a SizeTreeEntry) Less(b btree.Item) bool {
	return a.Size < b.(*SizeTreeEntry).Size
}


