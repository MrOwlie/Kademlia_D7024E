package kademlia

import (
	"os"
	"time"
	"fmt"
	"../d7024e"
	"../rpc"
)

const republishInterval time.Duration = 24 * time.Hour

func (kademlia *Kademlia) scheduleMessageBufferListGarbageCollect() {
	mbList := kademlia.MBList
	ticker := time.NewTicker(10 * time.Second)

	go func() {
		for {
			select {
			case <-ticker.C:
				fmt.Println("Collecting garbage")
				list := mbList.GarbageCollect()
				for i := 0; i < len(list); i++ {
					list[i].MessageChannel <- &rpc.Message{RpcType: rpc.TIME_OUT, RpcId: *list[i].RPCID, SenderId: *d7024e.NewRandomKademliaID(), RpcData: nil}
				}
			}
		}
	}()
}

func (kademlia *Kademlia) scheduleIdleBucketReExploration() {
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

func (kademlia *Kademlia) IdleBucketReExploration() {
	fmt.Println("Re-exploring buckets.")
	rTable := kademlia.routingTable
	var kademliaIDs []*d7024e.KademliaID = rTable.GetRefreshIDs()
	for i := 0; i < len(kademliaIDs); i++ {
		go kademlia.lookupProcedure(procedureContacts, kademliaIDs[i])
	}
}

func (kademlia *Kademlia) scheduleFileRepublish() {
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

func (kademlia *Kademlia) republishFiles() {
	fmt.Println("Republishing files...")
	metaData := kademlia.MetaData
	fileHashes := metaData.FilesToRepublish(republishInterval)
	nFiles := 0;
	for _, hash := range fileHashes {
		nFiles++;
		kademliaHash := d7024e.NewKademliaID(hash)
		closest := kademlia.LookupContact(kademliaHash)
		for _, contact := range closest {
			kademlia.sendStoreMessage(&contact, d7024e.NewRandomKademliaID(), kademliaHash, rpc.SENDER)
		}
	}
	fmt.Printf("%v files republished.\n", nFiles)
}

func (kademlia *Kademlia) scheduleCacheExpiredFileDeletion() {
	metaData := kademlia.MetaData
	ticker := time.NewTicker(60 * time.Minute)

	go func() {
		//kademlia := GetInstance()
		for {
			select {
			case <-ticker.C:
				fmt.Println("Deleting expired files...")
				filePaths := metaData.FilesToDelete()
				for i := 0; i < len(filePaths); i++ {
					os.Remove(filePaths[i])
				}
				fmt.Printf("%v files deleted.", len(filePaths))
			}
		}
	}()
}
