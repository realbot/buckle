package fileutils

import (
	"io/ioutil"
	"syscall"
	"testing"
)

func TestHashOf(t *testing.T) {
	const expectedHash = "8db963a7cac33aa7505af578d76cf0f5"

	f, err := ioutil.TempFile("", "samplehash")
	if err != nil {
		panic(err)
	}
	defer syscall.Unlink(f.Name())
	ioutil.WriteFile(f.Name(), []byte("Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nunc sollicitudin hendrerit dolor, at faucibus augue rutrum at.\n"), 0644)

	hash, err := HashOf(f.Name())
	if err != nil {
		t.Errorf("Unexpected hash error: %v", err)
	} else if hash != expectedHash {
		t.Errorf("Hash error expected %s actual %s", expectedHash, hash)
	}
}
