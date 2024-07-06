package filechecksum_test

import (
	// "fmt"
	"testing"

	"github.com/jessegalley/fsdupscan/internal/filechecksum"
	"github.com/stretchr/testify/assert"
)


func TestFileChecksum(t *testing.T) {
  assert.Equal(t, 1,1)
  file := "../../testdata/file4_9K.bin"
  checksum, err := filechecksum.CalculateChecksum(file)
  if err != nil {
    t.Fatalf("checksum calc failed on %v with %v", file, err)
  }
  // fmt.Println(checksum)
  assert.Equal(t, checksum, "1f9d7ebdff2be849d10147e7619245c882f94a39b60ebf5f798b9e3c5154a5a7")
}
