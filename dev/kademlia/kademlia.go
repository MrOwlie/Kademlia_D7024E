package kademlia

import (
	"encoding/json"
	"fmt"
	"sync"

	"time"

	"../d7024e"
	"../messageBufferList"
	"../network"
	"../routingTable"
	"../rpc"
)

type kademlia struct {
}

type kademliaMessage struct {
	returnType int
	contacts   []d7024e.Contact
}

var instance *kademlia
var once sync.Once
var selfContact d7024e.Contact = d7024e.NewContact(d7024e.NewRandomKademliaID(), "localhost")

const alpha int = 3
const valueK int = 20

const procedureContacts = 1
const procedureValue = 2
const returnContacts = 3
const returnHasValue = 4

func GetInstance() *kademlia {
	once.Do(func() {
		instance = &kademlia{}
		fmt.Println("Starting timed jobs")
		scheduleMessageBufferListGarbageCollect()
		fmt.Println("RPC timeout garbage collection started!")
		scheduleIdleBucketReExploration()
		fmt.Println("Idle bucket re-exploration started!")
		scheduleFileRepublish()
		fmt.Println("File republishing started!")
		scheduleCacheExpiredFileDeletion()
		fmt.Println("Cache expired file deletion started!")
		fmt.Println("All timed jobs started successfully!")
	})
	return instance
}

func (kademlia *kademlia) lookupProcedure(procedureType int, target *d7024e.KademliaID) {

	rTable := routingTable.GetInstance()
	candids := &d7024e.ContactCandidates{}
	queriedAddresses := make(map[string]bool)

	//select own k known closest node for first iteration
	candids.Append(rTable.FindClosestContacts(target, valueK))
	candids.Sort()

	var chans []chan kademliaMessage

	queriedAll := false
	for queriedAll == false {
		//get k closest or all candidates if lower
		nrQueries := valueK
		if candids.Len() < nrQueries {
			nrQueries = candids.Len()
		}
		iterCandids := candids.GetContacts(nrQueries)

		//check so the closest has not been queried, aka last iteration gave closer nodes
		if _, value := queriedAddresses[iterCandids[0].Address]; !value {
			//query the first alpha not queried
			queried := 0
			for _, contact := range iterCandids {
				if queried < alpha {
					if _, value := queriedAddresses[contact.Address]; !value {
						ch := make(chan kademliaMessage)
						chans = append(chans, ch)
						go kademlia.lookupSubProcedure(contact, target, procedureType, ch)
						queried = queried + 1
						queriedAddresses[contact.Address] = true
					}
				} else {
					break
				}
			}

			//wait for the queries to return data
			startTime := time.Now()
			incrementalLimit := 5
			startIndex := 0
			if len(chans) > alpha {
				startIndex = len(chans) - alpha
			}
			for {
				timeWaited := time.Since(startTime)
				allowedTimeouts := int(timeWaited.Seconds()) / incrementalLimit

				for i := startIndex; i < len(chans); i++ {
					select {
					case x, ok := <-chans[i]:
						if ok {
							if x.returnType == returnContacts {
								candids.Append(x.contacts)
							} else if x.returnType == returnHasValue {
								go kademlia.requestFile(target, x.contacts[0])
								return
							}
						}
					default:
						allowedTimeouts--
					}
				}

				if allowedTimeouts < 0 {
					time.Sleep(1000 * time.Millisecond)
				} else {
					break
				}
			}

			//go trhough all channels, save responses and remove clsoed channels
			for i := len(chans) - 1; i >= 0; i-- {
				select {
				case x, ok := <-chans[i]:
					if ok {
						if x.returnType == returnContacts {
							candids.Append(x.contacts)
						} else if x.returnType == returnHasValue {
							go kademlia.requestFile(target, x.contacts[0])
							return
						}
					}
					chans = append(chans[:i], chans[i+1:]...)
				default:
					continue
				}
			}

			//save k closest distinct contacts for next iteration
			candids.Sort()
			distinctContacts := candids.GetDistinctContacts(valueK)
			candids := &d7024e.ContactCandidates{}
			candids.Append(distinctContacts)

		} else { // last iteration gave no closer nodes, query all
			queriedAll = true

			//choose k closest nodes
			nrQueries := valueK
			if candids.Len() < nrQueries {
				nrQueries = candids.Len()
			}
			iterCandids := candids.GetContacts(nrQueries)

			//query all not already queried
			for _, contact := range iterCandids {
				if _, value := queriedAddresses[contact.Address]; !value {
					ch := make(chan kademliaMessage)
					chans = append(chans, ch)
					go kademlia.lookupSubProcedure(contact, target, procedureType, ch)
				}
			}

			//wait for responses or time-outs for all open queries
			for {
				if len(chans) == 0 {
					break
				}
				for i := len(chans) - 1; i >= 0; i-- {
					select {
					case x, ok := <-chans[i]:
						if ok {
							if x.returnType == returnContacts {
								candids.Append(x.contacts)
							} else if x.returnType == returnHasValue {
								go kademlia.requestFile(target, x.contacts[0])
								return
							}
						}
						chans = append(chans[:i], chans[i+1:]...)
					default:
						continue
					}
				}
				time.Sleep(1000 * time.Millisecond)
			}
		}

	}

}

func (kademlia *kademlia) lookupSubProcedure(target d7024e.Contact, toFind *d7024e.KademliaID, lookupType int, ch chan<- kademliaMessage) {
	//create and add its own messageBuffer to the singleton list
	net := network.GetInstance()
	rpcID := d7024e.NewRandomKademliaID()
	mBuffer := messageBufferList.NewMessageBuffer(rpcID)
	mBufferList := messageBufferList.GetInstance()
	mBufferList.AddMessageBuffer(mBuffer)

	//send different messages depending on type
	if lookupType == procedureContacts {
		net.SendFindContactMessage(&target, toFind, rpcID)
	} else if lookupType == procedureValue {
		net.SendFindDataMessage(&target, toFind, rpcID)
	}

	//wait until a response is retrieved
	mBuffer.WaitForResponse()
	message := mBuffer.ExtractMessage()

	//Return different flags and payload depending on file is found or contacts is returned
	if message.RpcType == rpc.CLOSEST_NODES {
		var contacts []d7024e.Contact
		json.Unmarshal(message.RpcData, contacts)
		retMessage := kademliaMessage{returnContacts, contacts}
		ch <- retMessage
		close(ch)
	} else if message.RpcType == rpc.HAS_VALUE {
		//return target so main routine knows which contact has the file
		contacts := []d7024e.Contact{target}
		retMessage := kademliaMessage{returnHasValue, contacts}
		ch <- retMessage
		close(ch)
	}
}

func (kademlia *kademlia) requestFile(file *d7024e.KademliaID, target d7024e.Contact) {

}

func (kademlia *kademlia) LookupContact(target *d7024e.KademliaID) {
	kademlia.lookupProcedure(procedureContacts, target)
}

func (kademlia *kademlia) LookupData(id string) {
	fileHash := d7024e.NewKademliaID(id)
	kademlia.lookupProcedure(procedureValue, fileHash)
}

func (kademlia *kademlia) Store(data []byte) {
	// TODO
}

func (kademlia *kademlia) Join(ip string, port int) {
	fmt.Printf("joining %q on port %d", ip, port)
	net := network.GetInstance()
	rpcID := d7024e.NewRandomKademliaID()
	bootstrapContact := d7024e.NewContact(d7024e.NewRandomKademliaID(), fmt.Sprintf("%s:%d", ip, port))
	mBuffer := messageBufferList.NewMessageBuffer(rpcID)
	mBufferList := messageBufferList.GetInstance()
	mBufferList.AddMessageBuffer(mBuffer)

	net.SendFindContactMessage(&bootstrapContact, selfContact.ID, rpcID)

	//wait until a response is
	fmt.Println("sent join message")
	mBuffer.WaitForResponse()
	message := mBuffer.ExtractMessage()
	fmt.Println("got response")

	var contacts []d7024e.Contact
	json.Unmarshal(message.RpcData, contacts)
	fmt.Println(contacts)
	for _, contact := range contacts {
		routingTable.GetInstance().AddContact(contact)
	}

}

func (kademlia *kademlia) ReturnLookupContact(target *d7024e.Contact) {
	// TODO
}

func (kademlia *kademlia) ReturnLookupData(hash string) {
	// TODO
}
