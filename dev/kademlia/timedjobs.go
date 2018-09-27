package kademlia

import (
	"os"
	"time"

	"../d7024e"
	"../messageBufferList"
	"../metadata"
	"../routingTable"
)

func scheduleMessageBufferListGarbageCollect() {
	mbList := messageBufferList.GetInstance()
	ticker := time.NewTicker(10 * time.Second)

	go func() {
		for {
			select {
			case <-ticker.C:
				mbList.GarbageCollect()
			}
		}
	}()
}

func scheduleIdleBucketReExploration() {
	ticker := time.NewTicker(60 * time.Minute)

	go func() {
		for {
			select {
			case <-ticker.C:
				kademlia := GetInstance()
				kademlia.IdleBucketReExploration()
			}
		}
	}()
}

func (kademlia *kademlia) IdleBucketReExploration() {
	rTable := routingTable.GetInstance()
	var kademliaIDs []*d7024e.KademliaID = rTable.GetRefreshIDs()

	for i := 0; i < len(kademliaIDs); i++ {
		go kademlia.lookupProcedure(procedureContacts, kademliaIDs[i])
	}
}

func scheduleFileRepublish() {
	metaData := metadata.GetInstance()
	ticker := time.NewTicker(60 * time.Minute)

	go func() {
		kademlia := GetInstance()
		for {
			select {
			case <-ticker.C:
				fileHashes := metaData.FilesToRepublish()

				for i := 0; i < len(fileHashes); i++ {
					go kademlia.StoreFile(fileHashes[i])
				}
			}
		}
	}()
}

func scheduleCacheExpiredFileDeletion() {
	metaData := metadata.GetInstance()
	ticker := time.NewTicker(60 * time.Minute)

	go func() {
		//kademlia := GetInstance()
		for {
			select {
			case <-ticker.C:
				filePaths := metaData.FilesToDelete()

				for i := 0; i < len(filePaths); i++ {
					os.Remove(filePaths[i])
				}
			}
		}
	}()
}
