package data

import (
    "testing"
    "os"
)

func TestSaveAndReadBuckleData(t *testing.T) {
    dataf := createTempBuckleDataFile()
    defer os.Remove(dataf.Name())
    WriteBuckleData(dataf, "foo", "8f24409843c176fa2c0b4690bfc94d15")
    dataf.Close()

    buckleData, err := ReadBuckleData(dataf.Name(), false)

    if err != nil {
        t.Errorf("Error reading buckle data file: %v", err)
    }
    if buckleData.Hashes["foo"] != "8f24409843c176fa2c0b4690bfc94d15" {
        t.Errorf("Data mismatch: expected 8f24409843c176fa2c0b4690bfc94d15 was %s", buckleData.Hashes["foo"])
    }
}

func TestResetReadBuckleData(t *testing.T) {
    dataf := createTempBuckleDataFile()
    defer os.Remove(dataf.Name())
    WriteBuckleData(dataf, "foo", "8f24409843c176fa2c0b4690bfc94d15")
    dataf.Close()

    buckleData, err := ReadBuckleData(dataf.Name(), true)

    if err != nil {
        t.Errorf("Error reading buckle data file: %v", err)
    }
    if len(buckleData.Hashes) != 0  {
        t.Error("Reset problem, hash not empty")
    }
}

func TestUpdateBuckleData(t *testing.T) {
    dataf := createTempBuckleDataFile()
    defer os.Remove(dataf.Name())
    WriteBuckleData(dataf, "foo", "8f24409843c176fa2c0b4690bfc94d15")
    dataf.Close()

    buckleData, err := ReadBuckleData(dataf.Name(), false)
    if err != nil {
        t.Errorf("Error reading buckle data file: %v", err)
    }

    buckleData.Hashes["foo"] = "12345678901234567890123456789012"
    buckleData.UpdateBuckleData(dataf.Name())
    reloadedData, err := ReadBuckleData(dataf.Name(), false)

    if err != nil {
        t.Errorf("Error updating buckle data file: %v", err)
    }
    if reloadedData.Hashes["foo"] != "12345678901234567890123456789012" {
        t.Errorf("Data was not updated: expected 12345678901234567890123456789012 was %s", reloadedData.Hashes["foo"])
    }
}

func TestChangedHashes(t *testing.T) {
    old := NewBuckleData()
    old.Hashes["change"] = "12345678901234567890123456789012"
    old.Hashes["delete"] = "99999999999999999999999999999999"

    updated := NewBuckleData()
    updated.Hashes["add"] = "00000000000000000000000000000000"
    updated.Hashes["change"] = "11111111111111111111111111111111"

    updated.CalculateChangedHashes(old)
    
    if updated.Hashes["change"] != "11111111111111111111111111111111" {
        t.Errorf("Data was not updated: expected 11111111111111111111111111111111 was %s", updated.Hashes["change"])
    }
    if updated.Hashes["add"] != "00000000000000000000000000000000" {
        t.Errorf("Data was not inserted: expected 00000000000000000000000000000000 was %s", updated.Hashes["change"])
    }
    if _, present := updated.Hashes["delete"]; present {
        t.Errorf("Data was not deleted")
    }
}
