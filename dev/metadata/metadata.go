package metadata

import (
	"fmt"
	"sync"
	"time"
)

const republishTime = 60

type MetaData struct {
	filePath        string
	hash            string
	isPinned        bool
	lastRepublish   time.Time
	timeOfInsertion time.Time
	timeToLive      time.Duration
}

type FileMetaData struct {
	fileData map[string]MetaData
	mutex    sync.Mutex
}

var instance FileMetaData
var once sync.Once

func GetInstance() *FileMetaData {
	once.Do(func() {
		instance = FileMetaData{}
		instance.fileData = make(map[string]MetaData)
	})
	return &instance
}

func (metaData *MetaData) TimeToRepublish() bool {
	elapsed := time.Since(metaData.lastRepublish)
	return (elapsed.Minutes() > republishTime)
}

func (metaData *MetaData) HasExpired() bool {
	timeAlive := metaData.timeOfInsertion.Add(metaData.timeToLive)
	currentTime := time.Now()
	return currentTime.After(timeAlive)
}

func (fileMetaData *FileMetaData) FilesToDelete() (filesToDelete []string) {
	fileMetaData.mutex.Lock()
	var keys []string

	for k := range fileMetaData.fileData {
		metaData := fileMetaData.fileData[k]
		if !metaData.isPinned && metaData.HasExpired() {
			keys = append(keys, k)
			filesToDelete = append(filesToDelete, metaData.filePath)
		}
	}

	for i := 0; i < len(keys); i++ {
		delete(fileMetaData.fileData, keys[i])
	}

	fileMetaData.mutex.Unlock()
	return
}

func (fileMetaData *FileMetaData) FilesToRepublish() (filesToRepublish []string) {
	fileMetaData.mutex.Lock()
	for k := range fileMetaData.fileData {
		metaData := fileMetaData.fileData[k]
		if !metaData.isPinned && metaData.HasExpired() {
			filesToRepublish = append(filesToRepublish, metaData.filePath)
		}
	}
	fileMetaData.mutex.Unlock()
	return
}

func (fileMetaData *FileMetaData) HasFile(hash string) bool {
	fileMetaData.mutex.Lock()
	_, has := fileMetaData.fileData[hash]
	fileMetaData.mutex.Unlock()
	return has
}

func (fileMetaData *FileMetaData) GetMetaData(hash string) (metaData *MetaData, found bool) {
	fileMetaData.mutex.Lock()
	*metaData, found = fileMetaData.fileData[hash]
	fileMetaData.mutex.Unlock()
	return
}

func (fileMetaData *FileMetaData) Pin(hash string) {
	fileMetaData.mutex.Lock()
	metaData, found := fileMetaData.fileData[hash]
	if found {
		metaData.isPinned = true
	} else {
		fmt.Println("No file with hash: %s \nwas found", hash)
	}
	fileMetaData.mutex.Unlock()
}

func (fileMetaData *FileMetaData) Unpin(hash string) {
	fileMetaData.mutex.Lock()
	metaData, found := fileMetaData.fileData[hash]
	if found {
		metaData.isPinned = false
		metaData.timeOfInsertion = time.Now()
	} else {
		fmt.Println("No file with hash: %s \nwas found", hash)
	}
	fileMetaData.mutex.Unlock()
}

func (fileMetaData *FileMetaData) AddFile(filePath string, hash string, pinned bool, timeToLive time.Duration) {
	if !fileMetaData.HasFile(hash) {
		fileMetaData.mutex.Lock()
		fileMetaData.fileData[hash] = MetaData{filePath, hash, pinned, time.Now(), time.Now(), timeToLive}
		fileMetaData.mutex.Unlock()
	} else {
		fmt.Println("File is already stored")
	}

}
