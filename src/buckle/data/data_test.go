package data

import (
	"testing"
)

func TestBoh(t *testing.T) {
	dataf := createTempBuckleDataFile()
	WriteBuckleData(dataf, "foo", "8f24409843c176fa2c0b4690bfc94d15")
	dataf.Close()
	
	data, err := ReadBuckleData(dataf.Name())
	if err != nil {
		t.Errorf("Error reading buckle data file: %v", err)
	}
	if data.Hashes["foo"] != "8f24409843c176fa2c0b4690bfc94d15" {
		t.Errorf("Data mismatch: expected bar was %s", data.Hashes["foo"])
	}
}