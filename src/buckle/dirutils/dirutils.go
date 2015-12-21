package dirutils

import (
	"crypto/md5"
	"fmt"
	"io"
	"math"
	"os"
	"path/filepath"
)

func ListFilesIn(path string) ([]string, error) {
	files := make([]string, 0, 100)

	visit := func(path string, f os.FileInfo, err error) error {
		var result error
		if err == nil {
			if !f.IsDir() {
				files = append(files, path)
			}
		} else {
			result = err
		}
		return result
	}

	err := filepath.Walk(path, visit)
	return files, err
}

const filechunk = 8192

func CalcHashOf(path string) (string, error) {
	var result string
	file, err := os.Open(path)
	if err == nil {
		defer file.Close()
		info, _ := file.Stat()
		filesize := info.Size()

		blocks := uint64(math.Ceil(float64(filesize) / float64(filechunk)))
		hash := md5.New()

		for i := uint64(0); i < blocks; i++ {
			blocksize := int(math.Min(filechunk, float64(filesize-int64(i*filechunk))))
			buf := make([]byte, blocksize)

			file.Read(buf)
			io.WriteString(hash, string(buf))
		}
		result = fmt.Sprintf("%x", hash.Sum(nil))
	}

	return result, err
}