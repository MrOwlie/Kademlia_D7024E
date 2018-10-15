package metadata

import (
	"fmt"
	"testing"
	"time"
)

func TestFileRepublish(t *testing.T) {
	expectedResult := make(map[string]bool)
	expectedResult["FFFFFFFFF0000000000000000000000000000001"] = true
	expectedResult["FFFFFFFFF0000000000000000000000000000002"] = true
	expectedResult["FFFFFFFFF0000000000000000000000000000003"] = true

	metadata := GetInstance()

	metadata.AddFile("filepath", "FFFFFFFFF0000000000000000000000000000001", true, time.Hour)
	metadata.AddFile("filepath", "FFFFFFFFF0000000000000000000000000000002", true, time.Hour)
	metadata.AddFile("filepath", "FFFFFFFFF0000000000000000000000000000003", true, time.Hour)

	time.Sleep(time.Second)

	metadata.AddFile("filepath", "FFFFFFFFF0000000000000000000000000000004", true, time.Hour)
	metadata.AddFile("filepath", "FFFFFFFFF0000000000000000000000000000005", true, time.Hour)
	metadata.AddFile("filepath", "FFFFFFFFF0000000000000000000000000000006", true, time.Hour)

	hashes := metadata.FilesToRepublish(time.Second)

	if len(hashes) != 3 {
		panic(fmt.Sprintf("Expected 3 files to republish got : %v\n", len(hashes)))
	}

	for _, v := range hashes {
		if expectedResult[v] {
			expectedResult[v] = false
		} else {
			panic(fmt.Sprintf("%v != %s", expectedResult, hashes))
		}
	}
}

func TestFileDeletion(t *testing.T) {
	expectedResult := make(map[string]bool)
	expectedResult["filepath1"] = true
	expectedResult["filepath2"] = true
	expectedResult["filepath3"] = true

	metadata := GetInstance()
	metadata.fileData = make(map[string]*MetaData)

	metadata.AddFile("filepath1", "FFFFFFFFF0000000000000000000000000000001", false, time.Millisecond)
	metadata.AddFile("filepath2", "FFFFFFFFF0000000000000000000000000000002", false, time.Millisecond)
	metadata.AddFile("filepath3", "FFFFFFFFF0000000000000000000000000000003", false, time.Millisecond)

	time.Sleep(time.Millisecond)

	metadata.AddFile("filepath4", "FFFFFFFFF0000000000000000000000000000004", false, 2*time.Second)
	metadata.AddFile("filepath5", "FFFFFFFFF0000000000000000000000000000005", false, 2*time.Second)
	metadata.AddFile("filepath6", "FFFFFFFFF0000000000000000000000000000006", false, 2*time.Second)

	filePaths := metadata.FilesToDelete()

	if len(filePaths) != 3 {
		panic(fmt.Sprintf("Expected 3 files to delete got : %v\n", len(filePaths)))
	}

	for _, v := range filePaths {
		if expectedResult[v] {
			expectedResult[v] = false
		} else {
			panic(fmt.Sprintf("%v != %s", expectedResult, filePaths))
		}
	}
}

func TestFilePinning(t *testing.T) {
	expectedResult := make(map[string]bool)
	expectedResult["filepath2"] = true
	expectedResult["filepath3"] = true

	metadata := GetInstance()
	metadata.fileData = make(map[string]*MetaData)

	metadata.AddFile("filepath1", "FFFFFFFFF0000000000000000000000000000001", false, time.Millisecond)
	metadata.AddFile("filepath2", "FFFFFFFFF0000000000000000000000000000002", false, time.Millisecond)
	metadata.AddFile("filepath3", "FFFFFFFFF0000000000000000000000000000003", false, time.Millisecond)
	metadata.Pin("FFFFFFFFF0000000000000000000000000000001")

	time.Sleep(time.Millisecond)

	filePaths := metadata.FilesToDelete()

	if len(filePaths) != 2 {
		panic(fmt.Sprintf("Expected 2 files to delete got : %v\n", len(filePaths)))
	}

	for _, v := range filePaths {
		if expectedResult[v] {
			expectedResult[v] = false
		} else {
			panic(fmt.Sprintf("%v != %s", expectedResult, filePaths))
		}
	}
}

func TestFileUnpinning(t *testing.T) {
	expectedResult := make(map[string]bool)
	expectedResult["filepath1"] = true
	expectedResult["filepath2"] = true
	expectedResult["filepath3"] = true

	metadata := GetInstance()
	metadata.fileData = make(map[string]*MetaData)

	metadata.AddFile("filepath1", "FFFFFFFFF0000000000000000000000000000001", true, 7*time.Second)
	metadata.AddFile("filepath2", "FFFFFFFFF0000000000000000000000000000002", false, 7*time.Second)
	metadata.AddFile("filepath3", "FFFFFFFFF0000000000000000000000000000003", false, 7*time.Second)
	metadata.Unpin("FFFFFFFFF0000000000000000000000000000001")

	time.Sleep(7 * time.Second)

	filePaths := metadata.FilesToDelete()

	if len(filePaths) != 3 {
		panic(fmt.Sprintf("Expected 3 files to delete got : %v\n", len(filePaths)))
	}

	for _, v := range filePaths {
		if expectedResult[v] {
			expectedResult[v] = false
		} else {
			panic(fmt.Sprintf("%v != %s", expectedResult, filePaths))
		}
	}
}

func TestFileRefresh(t *testing.T) {
	expectedResult := make(map[string]bool)
	expectedResult["FFFFFFFFF0000000000000000000000000000002"] = true
	expectedResult["FFFFFFFFF0000000000000000000000000000003"] = true

	metadata := GetInstance()
	metadata.fileData = make(map[string]*MetaData)

	metadata.AddFile("filepath", "FFFFFFFFF0000000000000000000000000000001", true, time.Hour)
	metadata.AddFile("filepath", "FFFFFFFFF0000000000000000000000000000002", true, time.Hour)
	metadata.AddFile("filepath", "FFFFFFFFF0000000000000000000000000000003", true, time.Hour)

	time.Sleep(2 * time.Second)
	metadata.RefreshFile("FFFFFFFFF0000000000000000000000000000001")

	hashes := metadata.FilesToRepublish(2 * time.Second)

	if len(hashes) != 2 {
		panic(fmt.Sprintf("Expected 2 files to delete got : %v\n", len(hashes)))
	}

	for _, v := range hashes {
		if expectedResult[v] {
			expectedResult[v] = false
		} else {
			panic(fmt.Sprintf("%v != %s", expectedResult, hashes))
		}
	}
}
