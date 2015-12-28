package data

import (
	"buckle/utils"
	"fmt"
	"os"
	"os/user"
	"io/ioutil"
	"strings"
)

type BuckleData struct {
	Hashes map[string]string 
}

func NewBuckleData() BuckleData {
	var bd BuckleData
	bd.Hashes = make(map[string]string)
	return bd
}

func (bd *BuckleData) UpdateBuckleData(buckleDataFilename string) {
	buckleWorkingFile := CreateTempBuckleData()
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
	const hashLen = 32
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

func CreateTempBuckleData() (f *os.File) {
	f, err := ioutil.TempFile("", "buckle")
	utils.Check(err)
	return
}

func WriteBuckleData(f *os.File, path string, hash string) {
	f.WriteString(fmt.Sprintf("%s=%s\n", hash, path))
}

func buckleDataFileExists(dataFilename string) bool {
	_, err := os.Stat(dataFilename)
	return err == nil
}