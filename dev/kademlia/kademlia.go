package kademlia

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"os"
	"sync"
	"time"

	"../d7024e"
	"../messageBufferList"
	"../metadata"
	"../routingTable"
	"../rpc"
)

type NetworkHandler interface {
	SendMessage(string, *[]byte)
	FetchFile(string, string) error
}

type kademlia struct {
	network NetworkHandler
}

type kademliaMessage struct {
	returnType int
	contacts   []d7024e.Contact
}

var instance *kademlia
var once sync.Once

var storagePath string = "/kademlia/storage/"
var downLoadPath string = "/kademlia/dowloads/"

var maxTTL float64 = float64(24 * time.Hour)
var coefficentTTL float64 = maxTTL / math.Exp(160)

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

func (kademlia *kademlia) SetNetworkHandler(handler NetworkHandler) {
	GetInstance().network = handler
}

func (kademlia *kademlia) lookupProcedure(procedureType int, target *d7024e.KademliaID) (*d7024e.ContactCandidates, *d7024e.Contact, bool) {

	rTable := routingTable.GetInstance()
	candids := &d7024e.ContactCandidates{}
	queriedAddresses := make(map[string]bool)

	//select own k known closest node for first iteration
	candids.Append(rTable.FindClosestContacts(target, valueK))
	candids.Sort()

	//Variables used in find value
	var fileHost *d7024e.Contact
	fileWasFound := false

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
								for i, c := range x.contacts {
									if c.ID == rTable.Me.ID {
										x.contacts = append(x.contacts[0:i-1], x.contacts[i+1:len(x.contacts)]...)
									}
								}
								candids.Append(x.contacts)
							} else if x.returnType == returnHasValue {
								for i, c := range candids.Contacts {
									if c.ID == candids.Contacts[0].ID {
										candids.Contacts = append(candids.Contacts[0:i-1], candids.Contacts[i+1:len(candids.Contacts)]...)
									}
								}
								if !fileWasFound {
									fileHost = &x.contacts[0]
									fileWasFound = true
								}
							}
						}
					default:
						allowedTimeouts--
					}
				}

				if fileWasFound {
					return candids, fileHost, fileWasFound
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
							for i, c := range x.contacts {
								if c.ID == rTable.Me.ID {
									x.contacts = append(x.contacts[0:i-1], x.contacts[i+1:len(x.contacts)]...)
								}
							}
							candids.Append(x.contacts)
						} else if x.returnType == returnHasValue {
							for i, c := range candids.Contacts {
								if c.ID == x.contacts[0].ID {
									candids.Contacts = append(candids.Contacts[0:i-1], candids.Contacts[i+1:len(candids.Contacts)]...)
								}
							}
							if !fileWasFound {
								fileHost = &x.contacts[0]
								fileWasFound = true
							}
						}
					}
					chans = append(chans[:i], chans[i+1:]...)
				default:
					continue
				}
			}

			if fileWasFound {
				return candids, fileHost, fileWasFound
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
								for i, c := range x.contacts {
									if c.ID == rTable.Me.ID {
										x.contacts = append(x.contacts[0:i-1], x.contacts[i+1:len(x.contacts)]...)
									}
								}
								candids.Append(x.contacts)
							} else if x.returnType == returnHasValue {
								for i, c := range candids.Contacts {
									if c.ID == candids.Contacts[0].ID {
										candids.Contacts = append(candids.Contacts[0:i-1], candids.Contacts[i+1:len(candids.Contacts)]...)
									}
								}
								if !fileWasFound {
									fileHost = &x.contacts[0]
									fileWasFound = true
								}
							}
						}
						chans = append(chans[:i], chans[i+1:]...)
					default:
						continue
					}
				}

				if fileWasFound {
					return candids, fileHost, fileWasFound
				}

				time.Sleep(1000 * time.Millisecond)
			}
		}

	}
	return candids, fileHost, fileWasFound

}

func (kademlia *kademlia) lookupSubProcedure(target d7024e.Contact, toFind *d7024e.KademliaID, lookupType int, ch chan<- kademliaMessage) {
	//create and add its own messageBuffer to the singleton list
	rpcID := d7024e.NewRandomKademliaID()
	mBuffer := messageBufferList.NewMessageBuffer(rpcID)
	mBufferList := messageBufferList.GetInstance()
	mBufferList.AddMessageBuffer(mBuffer)

	//send different messages depending on type
	if lookupType == procedureContacts {
		kademlia.sendFindContactMessage(&target, toFind, rpcID)
	} else if lookupType == procedureValue {
		kademlia.sendFindDataMessage(&target, toFind, rpcID)
	}

	//wait until a response is retrieved
	message := <-mBuffer.MessageChannel

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

func (kademlia *kademlia) LookupContact(target *d7024e.KademliaID) (closest *d7024e.ContactCandidates) {
	closest, _, _ = kademlia.lookupProcedure(procedureContacts, target)
	return
}

func (kademlia *kademlia) LookupData(id string) (filePath string, closest *d7024e.ContactCandidates) {
	fileHash := d7024e.NewKademliaID(id)
	closest, fileHost, fileWasFound := kademlia.lookupProcedure(procedureValue, fileHash)
	if fileWasFound {
		url := fileHost.Address + "/storage/" + id
		filePath = downLoadPath + id
		kademlia.network.FetchFile(url, filePath)
		kademlia.sendStoreMessage(&closest.Contacts[0], d7024e.NewRandomKademliaID(), fileHash, fileHost.Address)
	}
	return
}

func (kademlia *kademlia) StoreFile(filePath string) {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println(err)
		return
	}

	h := sha1.New()
	_, err = io.Copy(h, file)
	if err != nil {
		fmt.Println(err)
		return
	}
	file.Close()

	hash := h.Sum(nil)
	newFileName := hex.EncodeToString(hash)
	newPath := storagePath + newFileName
	os.Rename(filePath, newPath)

	kademliaHash := d7024e.NewKademliaID(newFileName)
	metadata.GetInstance().AddFile(newPath, newFileName, true, calcTimeToLive(kademliaHash))

	closest, _, _ := kademlia.lookupProcedure(procedureContacts, kademliaHash)

	for _, c := range closest.Contacts {
		kademlia.sendStoreMessage(&c, d7024e.NewRandomKademliaID(), kademliaHash, rpc.SENDER)
	}
}

func (kademlia *kademlia) Join(ip string, port int) bool {
	fmt.Printf("joining %q on port %d", ip, port)
	rpcID := d7024e.NewRandomKademliaID()
	bootstrapContact := d7024e.NewContact(d7024e.NewRandomKademliaID(), fmt.Sprintf("%s:%d", ip, port))
	mBuffer := messageBufferList.NewMessageBuffer(rpcID)
	mBufferList := messageBufferList.GetInstance()
	mBufferList.AddMessageBuffer(mBuffer)
	rt := routingTable.GetInstance()

	retry := 0
	for retry < 3 {
		fmt.Printf("Trying to connect. Try number %d", retry)
		kademlia.sendFindContactMessage(&bootstrapContact, rt.Me.ID, rpcID)

		//wait until a response is
		message := <-mBuffer.MessageChannel

		if message.RpcType == rpc.CLOSEST_NODES {

			var contacts rpc.ClosestNodes = rpc.ClosestNodes{}
			json.Unmarshal(message.RpcData, &contacts)
			for _, contact := range contacts.Closest {
				kademlia.addContact(&contact)
			}
			return true
		}
		rpcID = d7024e.NewRandomKademliaID()
		retry++
	}
	fmt.Printf("Failed to join network after three tries. Please try to connect to another node.")
	return false

}

func (kademlia *kademlia) addContact(contact *d7024e.Contact) {

	rt := routingTable.GetInstance()

	bucket := rt.Buckets[rt.GetBucketIndex(contact.ID)]

	if bucket.IsFull() {
		kademlia := GetInstance()
		rpcID := d7024e.NewRandomKademliaID()
		pingContact := bucket.List.Front().Value.(d7024e.Contact)
		pingContactElement := bucket.List.Front()
		mBuffer := messageBufferList.NewMessageBuffer(rpcID)
		mBufferList := messageBufferList.GetInstance()
		mBufferList.AddMessageBuffer(mBuffer)

		kademlia.sendPingMessage(&pingContact, rpcID)

		fmt.Println("sent ping message from bucket")

		message := <-mBuffer.MessageChannel
		fmt.Println("ping executed")

		if message.RpcType == rpc.PONG {
			fmt.Println("ping responded successfully, node alive.")
			bucket.AddContact(pingContact)
		} else {
			fmt.Println("ping without response, node dead.")
			bucket.List.Remove(pingContactElement)
			bucket.AddContact(*contact)
		}
	} else {
		bucket.AddContact(*contact)
	}

}

func calcTimeToLive(fileId *d7024e.KademliaID) (ttl time.Duration) {
	exponent := float64(routingTable.GetInstance().GetBucketIndex(fileId))
	ttl = time.Duration(coefficentTTL * math.Exp(exponent))
	return
}
