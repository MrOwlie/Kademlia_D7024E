package metadata

import (
	"fmt"
	"sync"
	"time"
)

type MetaData struct {
	filePath        string
	isPinned        bool
	shouldRepublish bool
	lastRepublish   time.Time
	timeOfInsertion time.Time
	timeToLive      time.Duration
}

func (metaData *MetaData) TimeToRepublish(republishInterval time.Duration) bool {
	elapsed := time.Since(metaData.lastRepublish)
	return (elapsed > republishInterval)
}

func (metaData *MetaData) HasExpired() bool {
	timeAlive := metaData.timeOfInsertion.Add(metaData.timeToLive)
	currentTime := time.Now()
	return currentTime.After(timeAlive)
}

func (metaData *MetaData) Refresh() {
	metaData.shouldRepublish = false
	metaData.timeOfInsertion = time.Now()
}

type FileMetaData struct {
	fileData map[string]*MetaData
	mutex    sync.Mutex
}

func (fileMetaData *FileMetaData) FilesToDelete() (filesToDelete []string) {
	fileMetaData.mutex.Lock()
	defer fileMetaData.mutex.Unlock()

	var keys []string

	for k, metaData := range fileMetaData.fileData {
		if !metaData.isPinned && metaData.HasExpired() {
			keys = append(keys, k)
			filesToDelete = append(filesToDelete, metaData.filePath)
		}
	}

	for i := 0; i < len(keys); i++ {
		delete(fileMetaData.fileData, keys[i])
	}
	return
}

func (fileMetaData *FileMetaData) FilesToRepublish(republishInterval time.Duration) (filesToRepublish []string) {
	fileMetaData.mutex.Lock()
	defer fileMetaData.mutex.Unlock()

	for hash, metaData := range fileMetaData.fileData {
		if metaData.isPinned && metaData.TimeToRepublish(republishInterval) {
			if metaData.shouldRepublish {
				filesToRepublish = append(filesToRepublish, hash)
			} else {
				metaData.shouldRepublish = true
			}
			metaData.lastRepublish = time.Now()
		}
	}
	return
}

func (fileMetaData *FileMetaData) HasFile(hash string) (hasFile bool) {
	fileMetaData.mutex.Lock()
	defer fileMetaData.mutex.Unlock()

	_, hasFile = fileMetaData.fileData[hash]
	return
}

func (fileMetaData *FileMetaData) RefreshFile(hash string) {
	fileMetaData.mutex.Lock()
	defer fileMetaData.mutex.Unlock()

	metaData, hasFile := fileMetaData.fileData[hash]
	if hasFile {
		metaData.Refresh()
	}
}

func (fileMetaData *FileMetaData) Pin(hash string) (found bool) {
	fileMetaData.mutex.Lock()
	defer fileMetaData.mutex.Unlock()
	metaData, found := fileMetaData.fileData[hash]
	if found {
		metaData.isPinned = true
		metaData.lastRepublish = time.Now()
	} else {
		fmt.Printf("No file with hash: %s \nwas found\n", hash)
	}
	return
}

func (fileMetaData *FileMetaData) Unpin(hash string) (found bool) {
	fileMetaData.mutex.Lock()
	defer fileMetaData.mutex.Unlock()

	metaData, found := fileMetaData.fileData[hash]
	if found {
		metaData.isPinned = false
		metaData.timeOfInsertion = time.Now()
	} else {
		fmt.Printf("No file with hash: %s \nwas found\n", hash)
	}
	return
}

func (fileMetaData *FileMetaData) AddFile(filePath string, hash string, pinned bool, timeToLive time.Duration) {
	fileMetaData.mutex.Lock()
	defer fileMetaData.mutex.Unlock()

	_, found := fileMetaData.fileData[hash]
	if found {
		fmt.Println("File already exists")
	} else {
		fileMetaData.fileData[hash] = &MetaData{filePath, pinned, true, time.Now(), time.Now(), timeToLive}
	}
}

func NewFileMetaData() *FileMetaData {
	instance := &FileMetaData{}
	instance.fileData = make(map[string]*MetaData)
	return instance
}
