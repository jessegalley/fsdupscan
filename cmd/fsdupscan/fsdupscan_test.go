package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)


func TestWalk (t *testing.T) {
  assert.Equal(t, 1, 1)

}

func TestValidateStartingDirs(t *testing.T) {
  positionalArgs := []string{"../../testdata/"}
  ok, err := validateStartingDirs(positionalArgs)
  if err != nil {
    t.Fatalf("error reading testdata/ in validateStartingDirs: %v", err)
  }
  assert.Equal(t, ok, true)

  positionalArgs = append(positionalArgs, "../../doesnotexist/")
  ok, _ = validateStartingDirs(positionalArgs)
  assert.Equal(t, ok, false)


  positionalArgs = []string{"../../testdata/top_dir1/", "../../testdata/top_dir2/"}
  ok, err = validateStartingDirs(positionalArgs)
  if err != nil {
    t.Fatalf("error reading testdata/ in validateStartingDirs: %v", err)
  }
  assert.Equal(t, ok, true)
}
