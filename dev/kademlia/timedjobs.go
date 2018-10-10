package kademlia

import (
	"fmt"
	"os"
	"time"

	"../d7024e"
	"../rpc"
)

const republishInterval time.Duration = 24 * time.Hour

func (kademlia *kademlia) scheduleMessageBufferListGarbageCollect() {
	mbList := kademlia.MBList
	ticker := time.NewTicker(10 * time.Second)

	go func() {
		for {
			select {
			case <-ticker.C:
				list := mbList.GarbageCollect()
				for i := 0; i < len(list); i++ {
					list[i].MessageChannel <- &rpc.Message{RpcType: rpc.TIME_OUT, RpcId: *list[i].RPCID, SenderId: *d7024e.NewRandomKademliaID(), RpcData: nil}
				}
			}
		}
	}()
}

func (kademlia *kademlia) scheduleIdleBucketReExploration() {
	ticker := time.NewTicker(60 * time.Minute)

	go func() {
		for {
			select {
			case <-ticker.C:
				kademlia.IdleBucketReExploration()
			}
		}
	}()
}

func (kademlia *kademlia) IdleBucketReExploration() {
	rTable := kademlia.routingTable
	var kademliaIDs []*d7024e.KademliaID = rTable.GetRefreshIDs()
	for i := 0; i < len(kademliaIDs); i++ {
		go kademlia.lookupProcedure(procedureContacts, kademliaIDs[i])
	}
}

func (kademlia *kademlia) scheduleFileRepublish() {
	ticker := time.NewTicker(60 * time.Minute)

	go func() {
		for {
			select {
			case <-ticker.C:
				kademlia.republishFiles()
			}
		}
	}()
}

func (kademlia *kademlia) republishFiles() {
	metaData := kademlia.MetaData
	fileHashes := metaData.FilesToRepublish(republishInterval)
	fmt.Println("Repub")
	for _, hash := range fileHashes {
		kademliaHash := d7024e.NewKademliaID(hash)
		closest := kademlia.LookupContact(kademliaHash)
		for _, contact := range closest {
			kademlia.sendStoreMessage(&contact, d7024e.NewRandomKademliaID(), kademliaHash, rpc.SENDER)
		}
	}
}

func (kademlia *kademlia) scheduleCacheExpiredFileDeletion() {
	metaData := kademlia.MetaData
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
