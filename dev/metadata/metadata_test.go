package metadata

import (
	"fmt"
	"testing"
	"time"
)

func TestFileRepublish(t *testing.T){
	expectedResult := make(map[string]bool)
	expectedResult["FFFFFFFFF0000000000000000000000000000001"] = true
	expectedResult["FFFFFFFFF0000000000000000000000000000002"] = true
	expectedResult["FFFFFFFFF0000000000000000000000000000003"] = true

	metadata := GetInstance()

	metadata.AddFile("filepath", "FFFFFFFFF0000000000000000000000000000001", false, time.Hour)
	metadata.AddFile("filepath", "FFFFFFFFF0000000000000000000000000000002", false, time.Hour)
	metadata.AddFile("filepath", "FFFFFFFFF0000000000000000000000000000003", false, time.Hour)

	time.Sleep(7*time.Second)

	metadata.AddFile("filepath", "FFFFFFFFF0000000000000000000000000000004", false, time.Hour)
	metadata.AddFile("filepath", "FFFFFFFFF0000000000000000000000000000005", false, time.Hour)
	metadata.AddFile("filepath", "FFFFFFFFF0000000000000000000000000000006", false, time.Hour)

	hashes := metadata.FilesToRepublish(7*time.Second)

	for _, v := range hashes {
		if expectedResult[v]{
			expectedResult[v] = false
		} else {
			panic(fmt.Sprintf("%v != %s", expectedResult, hashes))
		}
	}
}

func TestFileDeletion(t *testing.T){
	expectedResult := make(map[string]bool)
	expectedResult["filepath1"] = true
	expectedResult["filepath2"] = true
	expectedResult["filepath3"] = true

	metadata := GetInstance()

	metadata.AddFile("filepath1", "FFFFFFFFF0000000000000000000000000000001", false, 7*time.Second)
	metadata.AddFile("filepath2", "FFFFFFFFF0000000000000000000000000000002", false, 7*time.Second)
	metadata.AddFile("filepath3", "FFFFFFFFF0000000000000000000000000000003", false, 7*time.Second)

	time.Sleep(7*time.Second)

	metadata.AddFile("filepath4", "FFFFFFFFF0000000000000000000000000000004", false, 7*time.Second)
	metadata.AddFile("filepath5", "FFFFFFFFF0000000000000000000000000000005", false, 7*time.Second)
	metadata.AddFile("filepath6", "FFFFFFFFF0000000000000000000000000000006", false, 7*time.Second)

	filePaths := metadata.FilesToDelete()


	for _, v := range filePaths {
		if expectedResult[v]{
			expectedResult[v] = false
		} else {
			panic(fmt.Sprintf("%v != %s", expectedResult, filePaths))
		}
	}
}

func TestFilePinning(t *testing.T){
	expectedResult := make(map[string]bool)
	expectedResult["filepath2"] = true
	expectedResult["filepath3"] = true

	metadata := GetInstance()

	metadata.AddFile("filepath1", "FFFFFFFFF0000000000000000000000000000001", false, 7*time.Second)
	metadata.AddFile("filepath2", "FFFFFFFFF0000000000000000000000000000002", false, 7*time.Second)
	metadata.AddFile("filepath3", "FFFFFFFFF0000000000000000000000000000003", false, 7*time.Second)
	metadata.Pin("FFFFFFFFF0000000000000000000000000000001")

	time.Sleep(7*time.Second)

	metadata.AddFile("filepath4", "FFFFFFFFF0000000000000000000000000000004", false, 7*time.Second)
	metadata.AddFile("filepath5", "FFFFFFFFF0000000000000000000000000000005", false, 7*time.Second)
	metadata.AddFile("filepath6", "FFFFFFFFF0000000000000000000000000000006", false, 7*time.Second)

	filePaths := metadata.FilesToDelete()

	for _, v := range filePaths {
		if expectedResult[v]{
			expectedResult[v] = false
		} else {
			panic(fmt.Sprintf("%v != %s", expectedResult, filePaths))
		}
	}
}

func TestFileUnpinning(t *testing.T){
	expectedResult := make(map[string]bool)
	expectedResult["filepath1"] = true
	expectedResult["filepath2"] = true
	expectedResult["filepath3"] = true

	metadata := GetInstance()

	metadata.AddFile("filepath1", "FFFFFFFFF0000000000000000000000000000001", true, 7*time.Second)
	metadata.AddFile("filepath2", "FFFFFFFFF0000000000000000000000000000002", false, 7*time.Second)
	metadata.AddFile("filepath3", "FFFFFFFFF0000000000000000000000000000003", false, 7*time.Second)
	metadata.Unpin("FFFFFFFFF0000000000000000000000000000001")
	
	time.Sleep(7*time.Second)

	metadata.AddFile("filepath4", "FFFFFFFFF0000000000000000000000000000004", false, 7*time.Second)
	metadata.AddFile("filepath5", "FFFFFFFFF0000000000000000000000000000005", false, 7*time.Second)
	metadata.AddFile("filepath6", "FFFFFFFFF0000000000000000000000000000006", false, 7*time.Second)

	filePaths := metadata.FilesToDelete()

	for _, v := range filePaths {
		if expectedResult[v]{
			expectedResult[v] = false
		} else {
			panic(fmt.Sprintf("%v != %s", expectedResult, filePaths))
		}
	}
}