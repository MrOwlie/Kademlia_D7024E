package kademlia

import (
	"time"
	"os"

	"../d7024e"
	"../messageBufferList"
	"../routingTable"
	"../metadata"
)

func scheduleMessageBufferListGarbageCollect() {
	mbList := messageBufferList.GetInstance()
	ticker := time.NewTicker(5 * time.Minute)

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
	rTable := routingTable.GetInstance()
	ticker := time.NewTicker(5 * time.Minute)

	go func() {
		kademlia := GetInstance()
		for {
			select {
			case <-ticker.C:
				var kademliaIDs []*d7024e.KademliaID = rTable.GetRefreshIDs()

				for i := 0; i < len(kademliaIDs); i++ {
					go kademlia.lookupProcedure(procedureContacts, kademliaIDs[i])
				}
			}
		}
	}()
}

func scheduleFileRepublish() {
	metaData := metadata.GetInstance()
	ticker := time.NewTicker(5 * time.Minute)

	go func() {
		kademlia := GetInstance()
		for {
			select {
			case <-ticker.C:
				fileHashes := metaData.FilesToRepublish()

				for i := 0; i < len(fileHashes); i++ {
					go kademlia.Store(fileHashes[i])
				}
			}
		}
	}()
}

func scheduleCacheExpiredFileDeletion() {
	metaData := metadata.GetInstance()
	ticker := time.NewTicker(5 * time.Minute)

	go func() {
		kademlia := GetInstance()
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
