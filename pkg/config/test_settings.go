package config

import (
	"io/ioutil"
	"testing"
)

func SetupTempDir(t *testing.T) string {
	dir, err := ioutil.TempDir("./", "test")
	if err != nil {
		t.Fatalf("Unable to create Temp Dir : %s", dir)
	}
	return dir
}
