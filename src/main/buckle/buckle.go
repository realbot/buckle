package main

import (
	"buckle/utils"
	"buckle/fileutils"
	"buckle/data"
	"flag"
	"fmt"
	"log"
	"runtime"
)

func main() {
	var nCPU = flag.Int("numcpu", runtime.NumCPU(), "Number of CPU used")
	flag.Parse()
	
	runtime.GOMAXPROCS(*nCPU)
	log.Println("Number of CPUs: ", *nCPU)

	log.Println("Loading current hashes...")
	buckleDataFilename := data.BuckleDataFilename()
	
	buckleData, err := data.ReadBuckleData(buckleDataFilename)
	utils.CheckErrorMsg("Error reading buckle data file: ", err)
			
	files, err := fileutils.ListFilesIn("/home/realbot/temp")
	utils.CheckErrorMsg("Error reading dir content: ", err)

	fileHashes := calculateHashFor(files)
	for _, each := range fileHashes.CalculateChangedFiles(buckleData) {
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
