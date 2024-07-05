package filechecksum

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
)

func CalculateChecksum(filePath string) (string, error) {

  file, err := os.Open(filePath)
  if err != nil {
    return "", err
  }
  defer file.Close()

  hasher := sha256.New()
  _, err = io.Copy(hasher, file)
  if err != nil {
    return "", err
  }

  checksum := hasher.Sum(nil)
  return fmt.Sprintf("%x", checksum), nil
}
