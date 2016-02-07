package data

import (
	"buckle/utils"
	"fmt"
	"os"
	"os/user"
	"io/ioutil"
	"strings"
)

const hashLen = 32
	
type BuckleData struct {
	Hashes map[string]string 
}

func NewBuckleData() BuckleData {
	var bd BuckleData
	bd.Hashes = make(map[string]string)
	return bd
}

func (bd *BuckleData) UpdateBuckleData(buckleDataFilename string) {
	buckleWorkingFile := createTempBuckleDataFile()
	for path, currentHash := range bd.Hashes {
		WriteBuckleData(buckleWorkingFile, path, currentHash)
	}
	buckleWorkingFile.Close()
	os.Rename(buckleWorkingFile.Name(), buckleDataFilename)
}

func (bd *BuckleData) CalculateChangedFiles(oldData BuckleData) []string {
	result := make([]string, 0, len(bd.Hashes))
	for path, currentHash := range bd.Hashes {
		oldHash, oldHashExists := oldData.Hashes[path]
		if !oldHashExists || (oldHashExists && oldHash != currentHash) {
			result = append(result, path)
		}
	}
	return result
}

func BuckleDataFilename() string {
	usr, err := user.Current()
	if err != nil {
		panic(err)
	}
	return usr.HomeDir + "/.buckle"
}

func ReadBuckleData(dataFilename string) (BuckleData, error) {
	var result = NewBuckleData()
	var err error
	if buckleDataFileExists(dataFilename) {
		content, err := ioutil.ReadFile(dataFilename)
		if err == nil {
			for _, s := range strings.Split(string(content), "\n") {
				if len(s) > 0 {
					hash := s[:hashLen]
					path := s[hashLen+1:]
					result.Hashes[path] = hash
				}
			}
		}
	}
	return result, err
}

func WriteBuckleData(f *os.File, path string, hash string) {
	if len(hash) != hashLen {
		panic("hash mismatch " + hash)
	}
	f.WriteString(fmt.Sprintf("%s=%s\n", hash, path))
}

func createTempBuckleDataFile() (f *os.File) {
	f, err := ioutil.TempFile("", "buckle")
	utils.CheckError(err)
	return
}

func buckleDataFileExists(dataFilename string) bool {
	_, err := os.Stat(dataFilename)
	return err == nil
}
