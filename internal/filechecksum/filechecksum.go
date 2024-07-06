package filechecksum

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
)

func CalculateChecksum(filePath string) (string, error) {

  // return "1f9d7ebdff2be849d10147e7619245c882f94a39b60ebf5f798b9e3c5154a5a7", nil

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
  // checksum := "1234567"
  return fmt.Sprintf("%x", checksum), nil
}
