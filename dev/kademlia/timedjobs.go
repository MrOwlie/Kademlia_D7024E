package kademlia

import (
	"time"

	"../d7024e"
	"../messageBufferList"
	"../routingTable"
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
	kademlia := GetInstance()
	ticker := time.NewTicker(5 * time.Minute)

	go func() {
		for {
			select {
			case <-ticker.C:
				var kademliaIDs *[]d7024e.KademliaID = rTable.getRefreshIDs()

				for i := 0; i < len(kademliaIDs); i++ {
					kademlia.lookupProcedure(procedureContacts, kademliaIDs[i])
				}
			}
		}
	}()
}

func scheduleFileRepublish() {

}
