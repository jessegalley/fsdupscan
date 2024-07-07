package filechecksum

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"os"

	// "github.com/davecgh/go-spew/spew"
	"github.com/kalafut/imohash"
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

func CalculateChecksumQuick(filePath string) (string, error) {
  
  hasher := imohash.New()
  checksum, err := hasher.SumFile(filePath)
  if err != nil {
    return "", err
  }
  return ConvertToBase64(checksum), nil
}

func ConvertToBase64(input [16]byte) string {
	byteSlice := input[:]
	base64String := base64.StdEncoding.EncodeToString(byteSlice)
	return base64String
}
