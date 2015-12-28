package main

import (
	"buckle/fileutils"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"runtime"
	"strings"
)

func main() {
	nCPU := runtime.NumCPU()
	runtime.GOMAXPROCS(nCPU)
	log.Println("Number of CPUs: ", nCPU)

	log.Println("Loading current list...")
	buckleDataFilename := buckleDataFilename()
	buckleData := readBuckleData(buckleDataFilename)

	files, err := fileutils.ListFilesIn("/home/realbot/temp")
	check(err)

	fileHashes := calculateHashFor(files)

	buckleWorkingFile, err := ioutil.TempFile("", "buckle")
	check(err)

	for path, currentHash := range fileHashes {
		buckleWorkingFile.WriteString(fmt.Sprintf("%s=%s\n", currentHash, path))
		oldHash, oldHashExists := buckleData[path]
		if !oldHashExists || (oldHashExists && oldHash != currentHash) {
			fmt.Println(path)
		}
	}

	buckleWorkingFile.Close()
	//TODO move temp in .buckle...
}

func buckleDataFileExists(dataFilename string) bool {
	_, err := os.Stat(dataFilename)
	return err == nil
}

func buckleDataFilename() string {
	usr, err := user.Current()
	if err != nil {
		panic(err)
	}
	return usr.HomeDir + "/.buckle"
}

func readBuckleData(dataFilename string) map[string]string {
	const hashLen = 32
	var result = make(map[string]string)
	if buckleDataFileExists(dataFilename) {
		content, err := ioutil.ReadFile(dataFilename)
		if err != nil {
			panic("Error reading data file " + err.Error())
		}
		for _, s := range strings.Split(string(content), "\n") {
			if len(s) > 0 {
				hash := s[:hashLen]
				path := s[hashLen+1:]
				result[path] = hash
			}
		}
	}
	return result
}

func calculateHashFor(paths []string) map[string]string {
	const maxFileOpened = 50
	type item struct {
		path string
		hash string
		err  error
	}
	var result = make(map[string]string)
	var tokens = make(chan struct{}, maxFileOpened)

	ch := make(chan item, len(paths))
	for _, each := range paths {
		go func(p string) {
			tokens <- struct{}{}
			defer func() { <-tokens }()
			var it item
			it.path = p
			it.hash, it.err = fileutils.HashOf(p)
			ch <- it
		}(each)
	}

	for range paths {
		it := <-ch
		if it.err != nil {
			log.Printf("Error calculating hash for %s: %v\n", it.path, it.err)
		} else {
			result[it.path] = it.hash
		}
	}
	return result
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
