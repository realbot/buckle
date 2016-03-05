package fileutils

import (
	"io/ioutil"
	"syscall"
	"testing"
)

func TestHashOf(t *testing.T) {
	const expectedHash = "8db963a7cac33aa7505af578d76cf0f5"

	f, err := ioutil.TempFile("", "bucklehash")
	if err != nil {
		t.Errorf("Setup Error: %v", err)
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

func TestListFilesIn(t *testing.T) {
	parentDir, err := ioutil.TempDir("", "buckledir")
	if err != nil {
		t.Errorf("Setup Error: %v", err)
	}

	toexcludeDir, err := ioutil.TempDir(parentDir, "buckledir")
	if err != nil {
		t.Errorf("Setup Error: %v", err)
	}

	parentFile, err := ioutil.TempFile(parentDir, "bucklefile")
	if err != nil {
		t.Errorf("Setup Error: %v", err)
	}

	toexcludeFile, err := ioutil.TempFile(toexcludeDir, "bucklefile")
	if err != nil {
		t.Errorf("Setup Error: %v", err)
	}

	defer syscall.Unlink(toexcludeFile.Name())
	defer syscall.Unlink(parentFile.Name())
	defer syscall.Rmdir(toexcludeDir)
	defer syscall.Rmdir(parentDir)

	var noexclude Paths
	paths, err := ListFilesIn(parentDir, &noexclude)

	if err != nil {
		t.Errorf("Unexpected ListFilesIn: %v", err)
	} else if len(paths) != 2 {
		t.Errorf("Expected two files, got %v", paths)
	}

	var withexclude Paths
	withexclude.Set(toexcludeDir)
	singlepath, err := ListFilesIn(parentDir, &withexclude)

	if err != nil {
		t.Errorf("Unexpected ListFilesIn: %v", err)
	} else if len(singlepath) != 1 && singlepath[0] != parentFile.Name() {
		t.Errorf("Expected just one files, got %v", singlepath)
	}
}
