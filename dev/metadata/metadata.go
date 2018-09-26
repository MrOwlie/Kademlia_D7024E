package metadata

import (
	"time"
	"sync"
	"fmt"
)



const republishTime = 60

type MetaData struct{
	filePath string
	hash string
	isPinned bool 
	lastRepublish time.Time
	timeOfInsertion time.Time
	timeToLive time.Duration
}


type FileMetaData struct{
	fileData []MetaData
	mutex    sync.Mutex
}

var instance FileMetaData
var once sync.Once

func GetInstance() *FileMetaData {
	once.Do(func() {	
		instance = FileMetaData{}

	})
	return &instance
}


func (metaData *MetaData) TimeToRepublish() bool{
	elapsed := time.Since(metaData.lastRepublish)
	return (elapsed.Minutes() > republishTime)
}

func (metaData *MetaData) HasExpired() bool{
	timeAlive := metaData.timeOfInsertion.Add(metaData.timeToLive)
	currentTime := time.Now()
	return currentTime.After(timeAlive)
}


func (fileMetaData *FileMetaData) FilesToDelete() (filesToDelete []string){
	fileMetaData.mutex.Lock()
	var indexes []int

	for i := 0; i < len(fileMetaData.fileData); i++{
		if !fileMetaData.fileData[i].isPinned && fileMetaData.fileData[i].HasExpired(){
			indexes = append(indexes, i) 
			filesToDelete = append(filesToDelete, fileMetaData.fileData[i].filePath) 
		}	
	}

	for i := 0; i < len(indexes); i++{
		fileMetaData.fileData = append(fileMetaData.fileData[:indexes[i]], fileMetaData.fileData[indexes[i]+1:]...)	
	}

	fileMetaData.mutex.Unlock()
	return 
}

func (fileMetaData *FileMetaData) FilesToRepublish() (filesToRepublish []string){
	fileMetaData.mutex.Lock()
	for i := 0; i < len(fileMetaData.fileData); i++{
		if fileMetaData.fileData[i].TimeToRepublish(){
			filesToRepublish = append(filesToRepublish, fileMetaData.fileData[i].filePath)
			fileMetaData.fileData[i].lastRepublish = time.Now()
		}	
	}
	fileMetaData.mutex.Unlock()
	return
}

func (fileMetaData *FileMetaData) HasFile(hash string) bool{
	fileMetaData.mutex.Lock()
	for i := 0; i < len(fileMetaData.fileData); i++{
		if fileMetaData.fileData[i].hash == hash{
			fileMetaData.mutex.Unlock()
			return true 
		}	
	
	fileMetaData.mutex.Unlock()
	return false
}

func (fileMetaData *FileMetaData) getFileMetaData(hash string) (metaData *MetaData, found bool){
	for i := 0; i < len(fileMetaData.fileData); i++{
		if fileMetaData.fileData[i].hash == hash{
			found = true
		}	
	}
	return
}

func (fileMetaData *FileMetaData) Pin(hash string){
	fileMetaData.mutex.Lock()
	metaData, found := fileMetaData.getFileMetaData(hash)
	if found {
		metaData.isPinned = true
	} else {
		fmt.Println("No file with hash: %s \nwas found", hash)
	}
	fileMetaData.mutex.Unlock()
}

func (fileMetaData *FileMetaData) Unpin(hash string){
	fileMetaData.mutex.Lock()
	metaData, found := fileMetaData.getFileMetaData(hash)
	if found {
		metaData.isPinned = false
		metaData.timeOfInsertion = time.Now()
	} else {
		fmt.Println("No file with hash: %s \nwas found", hash)
	}
	fileMetaData.mutex.Unlock()
}

func (fileMetaData *FileMetaData) AddFile(filePath string, hash string, pinned bool, timeToLive time.Duration){
	fileMetaData.mutex.Lock()
	fileMetaData.fileData = append(fileMetaData.fileData, MetaData{filePath, hash, pinned, time.Now(), time.Now(), timeToLive})
	fileMetaData.mutex.Unlock()
}

