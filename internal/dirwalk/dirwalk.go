package dirwalk

import (
	"errors"
	"log/slog"
	"os"
	"path/filepath"
	"sync"
	"syscall"

	// "github.com/davecgh/go-spew/spew"
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
  fileCh := make(chan *WalkedFile, 1)
  for _, dir := range dirs {
    wg.Add(1)
    go WalkDir(dir, &wg, fileCh)
  }
  
  return fileCh, &wg
}

func WalkDir(dir string, wg *sync.WaitGroup, fileCh chan<- *WalkedFile) {
  defer wg.Done()

  visit := func (path string, file os.FileInfo, err error) error {
    if file.IsDir() && path != dir {
      wg.Add(1)
      go WalkDir(path, wg, fileCh)
      return filepath.SkipDir
    }
    if file.Mode().IsRegular() {
      // spew.Dump(file.Sys().(*syscall.Stat_t).Ino)
      inode := file.Sys().(*syscall.Stat_t).Ino
      fileCh <- NewWalkedFile(filepath.Join(path, file.Name()), inode, file.Size())   
      slog.Debug("WalkDir visit file found", "path", path, "name", file.Name(), "size", file.Size())
    }
    return nil
  }

  filepath.Walk(dir, visit)
}

// func Walk(dir string, fileCh chan<- os.DirEntry, followSym bool) {
//   entries, err := readDirRegular(dir)
//   if err != nil {
//     // we shouldn't get here
//     // panic("something broke when reading dir in walk()")
//     panic(err)
//   }
//
//   for _, entry := range entries {
//     if entry.IsDir() {
//       Walk(filepath.Join(dir, entry.Name()), fileCh, followSym)
//     } else if entry.Type().IsRegular() {
//       fileCh <- entry
//     } else if isSymlink(entry) {
//       if !followSym {
//         continue
//       }
//       // target, err := resolveSymlink(entry, dir)
//       // if err != nil {
//       //   panic(err)
//       // }
//       // 
//       // statTarget := target
//       // if !filepath.IsAbs(target) {
//       //   statTarget = filepath.Join(dir, target)
//       // }
//       // fileInfo, err := os.Stat(statTarget)
//       // if err != nil {
//       //   //TODO: what does a failed stat() mean here? broken symlink?
//       //   panic(err)
//       // }
//     } else {
//       slog.Debug("dirwalk::walk() unknown entry.Type()", "Type()", entry.Type())
//       continue
//     }
//   }
// }

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


