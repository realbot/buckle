package main

import (
	"buckle/fileutils"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"strings"
)

func main() {
	log.Println("Loading current list...")
	//buckleData := readBuckleData(buckleDataFilename())

	files, err := fileutils.ListFilesIn("/home/realbot/tmp1")
	if err != nil {
		panic(err)
	}

	filehashes := calculateHashFor(files)

	/*for i, each := range files {
		fmt.Printf("%d %s\n", i, each)
	}*/
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

func readBuckleData(dataFilename string) []string {
	var result []string
	if buckleDataFileExists(dataFilename) {
		content, err := ioutil.ReadFile(dataFilename)
		if err != nil {
			panic("Error reading data file " + err.Error())
		}
		result = strings.Split(string(content), "\n")
	}
	return result
}

func calculateHashFor(paths []string) map[string]string {
	type item struct {
		path string
		hash string
		err  error
	}
	var result map[string]string

	ch := make(chan item, len(paths))
	for _, each := range paths {
		go func(p string) {
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
		}
		//...to be continued... :-)
	}
	return result
}
