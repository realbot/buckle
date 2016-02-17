package data

import (
    "testing"
)

func TestSaveAndReadBuckleData(t *testing.T) {
    dataf := createTempBuckleDataFile()
    WriteBuckleData(dataf, "foo", "8f24409843c176fa2c0b4690bfc94d15")
    dataf.Close()

    buckleData, err := ReadBuckleData(dataf.Name())

    if err != nil {
        t.Errorf("Error reading buckle data file: %v", err)
    }
    if buckleData.Hashes["foo"] != "8f24409843c176fa2c0b4690bfc94d15" {
        t.Errorf("Data mismatch: expected 8f24409843c176fa2c0b4690bfc94d15 was %s", buckleData.Hashes["foo"])
    }
}

func TestUpdateBuckleData(t *testing.T)  {
    dataf := createTempBuckleDataFile()
    WriteBuckleData(dataf, "foo", "8f24409843c176fa2c0b4690bfc94d15")
    dataf.Close()
    
    buckleData, err := ReadBuckleData(dataf.Name())
    if err != nil {
        t.Errorf("Error reading buckle data file: %v", err)
    }

    buckleData.Hashes["foo"] = "12345678901234567890123456789012"
    buckleData.UpdateBuckleData(dataf.Name())
    reloadedData, err := ReadBuckleData(dataf.Name())

    if err != nil {
        t.Errorf("Error updating buckle data file: %v", err)
    }
    if reloadedData.Hashes["foo"] != "12345678901234567890123456789012" {
        t.Errorf("Data was not updated: expected 12345678901234567890123456789012 was %s", reloadedData.Hashes["foo"])
    }    
}
