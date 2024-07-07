package filechecksum_test

import (
	"testing"

	"github.com/jessegalley/fsdupscan/internal/filechecksum"
	"github.com/stretchr/testify/assert"
)


func TestFileChecksum(t *testing.T) {
  file := "../../testdata/file4_9K.bin"
  checksum, err := filechecksum.CalculateChecksum(file)
  if err != nil {
    t.Fatalf("checksum calc failed on %v with %v", file, err)
  }
  // fmt.Println(checksum)
  assert.Equal(t, checksum, "1f9d7ebdff2be849d10147e7619245c882f94a39b60ebf5f798b9e3c5154a5a7")
}

func TestFileChecksumQuick(t *testing.T) {
  file := "../../testdata/file4_9K.bin"
  checksum, err := filechecksum.CalculateChecksumQuick(file)
  if err != nil {
    t.Fatalf("checksum calc failed on %v with %v", file, err)
  }
  assert.Equal(t, checksum, "gEgF1tW+ca3eOTCKg/SjWg==")
}
