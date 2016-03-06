package main

import (
	"buckle/data"
	"buckle/fileutils"
	"buckle/utils"
	"flag"
	"fmt"
	"log"
	"runtime"
)

const VERSION = "0.1"

func main() {
	var nCPU = flag.Int("numcpu", runtime.NumCPU(), "Number of CPU used")
    var from = flag.String("from", utils.CurrentUser().HomeDir, "Starting path to process")
	var exclude fileutils.Paths
	flag.Var(&exclude, "exclude", "Directory to exclude")
	flag.Parse()

	runtime.GOMAXPROCS(*nCPU)
    log.Printf("Buckle (%v)\n", VERSION)
	log.Println("Number of CPUs: ", *nCPU)
    log.Println("Starting path: ", *from)
	for _, e := range exclude {
		log.Println("Excluded: ", e)
	}

	log.Println("Loading current hashes...")
	buckleDataFilename := data.BuckleDataFilename()

	buckleData, err := data.ReadBuckleData(buckleDataFilename)
	utils.CheckErrorMsg("Error reading buckle data file: ", err)

	files, err := fileutils.ListFilesIn(*from, &exclude)
	utils.CheckErrorMsg("Error reading dir content: ", err)

	fileHashes := calculateHashFor(files)
	for _, each := range fileHashes.CalculateChangedHashes(buckleData) {
		fmt.Println(each)
	}
	fileHashes.UpdateBuckleData(buckleDataFilename)
}

func calculateHashFor(paths []string) data.BuckleData {
	const maxFileOpened = 50
	type item struct {
		path string
		hash string
		err  error
	}
	var result = data.NewBuckleData()
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
			result.Hashes[it.path] = it.hash
		}
	}
	return result
}
