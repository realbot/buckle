package fileutils

import (
	"io/ioutil"
	"testing"
    "os"
    "path/filepath"
)

func TestHashOf(t *testing.T) {
	const expectedHash = "8db963a7cac33aa7505af578d76cf0f5"

	f := createTempFile("", t)
	defer os.Remove(f.Name())
	ioutil.WriteFile(f.Name(), []byte("Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nunc sollicitudin hendrerit dolor, at faucibus augue rutrum at.\n"), 0644)

	hash, err := HashOf(f.Name())
	if err != nil {
		t.Errorf("Unexpected hash error: %v", err)
	} else if hash != expectedHash {
		t.Errorf("Hash error expected %s actual %s", expectedHash, hash)
	}
}

func TestListFilesIn(t *testing.T) {
	parentDir := createTempDir("", t)
	parentFile := createTempFile(parentDir, t)
	toexcludeDir := createTempDir(parentDir, t)
	createTempFile(toexcludeDir, t)
    toexcludeSubDir := createTempDir(toexcludeDir, t)
    createTempFile(toexcludeSubDir, t)
	defer os.RemoveAll(parentDir)

	var noexclude Paths
	paths, err := ListFilesIn(parentDir, &noexclude)

	if err != nil {
		t.Errorf("Unexpected ListFilesIn: %v", err)
	} else if len(paths) != 3 {
		t.Errorf("Expected three files, got %v", paths)
	}

	var withexclude Paths
	withexclude.Set(toexcludeDir)
	singlepath, err := ListFilesIn(parentDir, &withexclude)

	if err != nil {
		t.Errorf("Unexpected ListFilesIn: %v", err)
	} else if len(singlepath) != 1 || singlepath[0] != parentFile.Name() {
		t.Errorf("Expected just one files, got %v", singlepath)
	}
}

func TestListFilesInSymLink(t *testing.T) {
    parentDir := createTempDir("", t)
    parentFile := createTempFile(parentDir, t)
    symlink := filepath.Join(parentDir, "symlink")
    os.Symlink(parentFile.Name(), symlink)
	defer os.RemoveAll(parentDir)

    var noexclude Paths
	paths, err := ListFilesIn(parentDir, &noexclude)

    if err != nil {
		t.Errorf("Unexpected ListFilesIn: %v", err)
    } else if len(paths) != 1 || paths[0] != parentFile.Name() {
		t.Errorf("Expected just one files, got %v", paths)
	}
}


func createTempDir(path string, t *testing.T) string {
 	someDir, err := ioutil.TempDir(path, "buckledir")
	if err != nil {
		t.Errorf("Setup Error: %v", err)
	}
    return someDir
}

func createTempFile(path string, t *testing.T) *os.File {
   	someFile, err := ioutil.TempFile(path, "bucklefile")
	if err != nil {
		t.Errorf("Setup Error: %v", err)
	}
    return someFile
}