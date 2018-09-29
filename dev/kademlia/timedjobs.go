package kademlia

import (
	"os"
	"time"

	"../d7024e"
	"../messageBufferList"
	"../metadata"
	"../routingTable"
	"../rpc"
)
const republishInterval time.Duration = 24*time.Hour

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
	ticker := time.NewTicker(60 * time.Minute)

	go func() {
		kademlia := GetInstance()
		for {
			select {
			case <-ticker.C:
				kademlia.republishFiles()
			}
		}
	}()
}

func (kademlia *kademlia) republishFiles() {
	metaData := metadata.GetInstance()
	fileHashes := metaData.FilesToRepublish(republishInterval)
	
	for _, hash := range fileHashes {
		kademliaHash := d7024e.NewKademliaID(hash)
		closest := kademlia.LookupContact(kademliaHash)
		for _, contact := range closest.Contacts {
			kademlia.sendStoreMessage(&contact, d7024e.NewRandomKademliaID(), kademliaHash, rpc.SENDER)
		}
	}
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
