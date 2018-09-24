package metadata

import (
	"time"
	"sync"
)



const republishTime = 60

type MetaData struct{
	FilePath string
	Hash string 
	lastRepublish time.Time
	expirationDate time.Time
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
	currentTime := time.Now()
	return currentTime.After(metaData.expirationDate)
}


func (fileMetaData *FileMetaData) HasFile(hash string) bool{
	fileMetaData.mutex.Lock()
	for i := 0; i < len(fileMetaData.fileData); i++{
		if fileMetaData.fileData[i].Hash == hash{
			fileMetaData.mutex.Unlock()
			return true 
		}	
	}
	fileMetaData.mutex.Unlock()
	return false
}

func (fileMetaData *FileMetaData) FilesToDelete() (filesToDelete []string){
	fileMetaData.mutex.Lock()
	var indexes []int

	for i := 0; i < len(fileMetaData.fileData); i++{
		if fileMetaData.fileData[i].HasExpired(){
			indexes = append(indexes, i) 
			filesToDelete = append(filesToDelete, fileMetaData.fileData[i].FilePath) 
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
			filesToRepublish = append(filesToRepublish, fileMetaData.fileData[i].FilePath)
			fileMetaData.fileData[i].lastRepublish = time.Now()
		}	
	}
	fileMetaData.mutex.Unlock()
	return
}

func (fileMetaData *FileMetaData) AddFile(filePath string, hash string, expirationDate time.Time){
	fileMetaData.mutex.Lock()
	fileMetaData.fileData = append(fileMetaData.fileData, MetaData{filePath, hash, time.Now(), expirationDate})
	fileMetaData.mutex.Unlock()
}

