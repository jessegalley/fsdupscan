package dirwalk

import (
	"errors"
	"log/slog"
	"os"
	"path/filepath"
	"sync"
	"syscall"
)

type WalkedFile struct {
  Path string 
  Inode uint64
  Size  int64 
}

func NewWalkedFile(path string, inode uint64, size int64) *WalkedFile {
  return &WalkedFile{
    Path: path,
    Inode: inode,
    Size: size,
  }
}

func Walk(dirs ...string) (<-chan *WalkedFile, *sync.WaitGroup) {
  var wg sync.WaitGroup
  fileCh := make(chan *WalkedFile, 100)
  for _, dir := range dirs {
    wg.Add(1)
    go WalkDir(dir, &wg, fileCh)
  }
  
  return fileCh, &wg
}

func WalkDir(dir string, wg *sync.WaitGroup, fileCh chan<- *WalkedFile) {
  defer wg.Done()

  visit := func (path string, file os.FileInfo, err error) error {
    // sometimes a file will get removed between the time it is listed  
    // and the time where it is to be Stat'd here. If this is the case the 
    // path will exist, but the actual FileInfo will be nil. so we wil 
    // return the visit function early to prevent segfaults.
    if file == nil {
      slog.Debug("visit func: file is nil (deleted during scan)", "path", path )
      return nil
    }

    // because of the sometimes nil FileInfo, to be safe we're going to return 
    // early if for some reason we have a garbage (empty) path as well
    if path == "" {
      slog.Debug("visit func: empty file path, skipping", "path", path )
      return nil
    }

    if file.IsDir() && path != dir {
      wg.Add(1)
      go WalkDir(path, wg, fileCh)
      return filepath.SkipDir
    }

    if file.Mode().IsRegular() {
      inode := file.Sys().(*syscall.Stat_t).Ino
      // fileCh <- NewWalkedFile(filepath.Join(path, file.Name()), inode, file.Size())   
      // path var inside the visit function contains the full path of files, not
      // just the parent dir
      fileCh <- NewWalkedFile(path, inode, file.Size())   
      // slog.Debug("WalkDir visit file found", "path", path, "name", file.Name(), "size", file.Size())
    }
    return nil
  }

  filepath.Walk(dir, visit)
}

// readDirRegular reads the entries in a directory path
func readDirRegular(dir string) ([]os.DirEntry, error) { 
  entries, err := os.ReadDir(dir)
  if err != nil {
    return nil, err
  }

  return entries, nil
}

// isSymlink chcks to see if the given file DirEntry is is a symlink. 
// returns bool 
func isSymlink(file os.DirEntry) bool {
  mode := file.Type()
  return mode & os.ModeSymlink != 0
}

// resolveSymlink attempts to resolve a the target of a symlink at File 
// located at path. will return an error if file is not a symlink or fails 
// to be read.  Returns a string of the link target which could be a 
// relative or absolute path.
func resolveSymlink(file os.DirEntry, dir string) (string, error) {
  linkpath := filepath.Join(dir, file.Name())
  if !isSymlink(file){
    slog.Debug("diwalk resolveSymlink() not a link", "link", linkpath)
    return "", errors.New("calling resolveSymlink on not a link")
  }

  target, err := os.Readlink(linkpath)
  if err != nil {
    slog.Debug("diwalk resolveSymlink() failed to reoslve", "link", linkpath)
    return "", err
  }

  return target, nil
}


