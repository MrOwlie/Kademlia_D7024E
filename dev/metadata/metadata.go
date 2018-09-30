package metadata

import (
	"fmt"
	"sync"
	"time"
)



type FileMetaData struct {
	fileData map[string]*MetaData
	mutex    sync.Mutex
}

var instance FileMetaData
var once sync.Once

func GetInstance() *FileMetaData {
	once.Do(func() {
		instance = FileMetaData{}
		instance.fileData = make(map[string]*MetaData)
	})
	return &instance
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
	metaData.lastRepublish = time.Now()
	metaData.timeOfInsertion = time.Now()
}


type MetaData struct {
	filePath        string
	isPinned        bool
	lastRepublish   time.Time
	timeOfInsertion time.Time
	timeToLive      time.Duration
}


func (fileMetaData *FileMetaData) FilesToDelete() (filesToDelete []string) {
	fileMetaData.mutex.Lock()
	defer fileMetaData.mutex.Unlock()

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
	return
}

func (fileMetaData *FileMetaData) FilesToRepublish(republishInterval time.Duration) (filesToRepublish []string) {
	fileMetaData.mutex.Lock()
	defer fileMetaData.mutex.Unlock()
	
	for hash := range fileMetaData.fileData {
		metaData := fileMetaData.fileData[hash]
		if !metaData.isPinned && metaData.TimeToRepublish(republishInterval) {
			filesToRepublish = append(filesToRepublish, hash)
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

	metadata, hasFile := fileMetaData.fileData[hash]
	if hasFile {
		metadata.Refresh()
	}
}

func (fileMetaData *FileMetaData) Pin(hash string) {
	fileMetaData.mutex.Lock()
	defer fileMetaData.mutex.Unlock()
	metaData, found := fileMetaData.fileData[hash]
	if found {
		metaData.isPinned = true
	} else {
		fmt.Printf("No file with hash: %s \nwas found\n", hash)
	}
}

func (fileMetaData *FileMetaData) Unpin(hash string) {
	fileMetaData.mutex.Lock()
	defer fileMetaData.mutex.Unlock()

	metaData, found := fileMetaData.fileData[hash]
	if found {
		metaData.isPinned = false
		metaData.timeOfInsertion = time.Now()
	} else {
		fmt.Printf("No file with hash: %s \nwas found\n", hash)
	}
}

func (fileMetaData *FileMetaData) AddFile(filePath string, hash string, pinned bool, timeToLive time.Duration) {
	fileMetaData.mutex.Lock()
	defer fileMetaData.mutex.Unlock()

	_, found := fileMetaData.fileData[hash]
	if found {
		fmt.Println("File already exists")
	} else {
		fileMetaData.fileData[hash] = &MetaData{filePath, pinned, time.Now(), time.Now(), timeToLive}
	}
}
