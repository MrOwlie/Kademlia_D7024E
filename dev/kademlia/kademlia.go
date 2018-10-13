package kademlia

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"os"
	"path/filepath"
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

type Kademlia struct {
	network      NetworkHandler
	routingTable *routingTable.RoutingTable
	MBList       *messageBufferList.MessageBufferList
	MetaData     *metadata.FileMetaData
}

type kademliaMessage struct {
	returnType int
	contacts   []d7024e.Contact
}

var storagePath string
var downLoadPath string

var maxTTL float64 = float64(24 * time.Hour)
var coefficentTTL float64 = maxTTL / math.Exp(160)

const alpha int = 3
const valueK int = 20

const procedureContacts = 1
const procedureValue = 2
const returnContacts = 3
const returnHasValue = 4
const returnTimedOut = 5

func NewKademliaObject(RT *routingTable.RoutingTable, MBL *messageBufferList.MessageBufferList, MD *metadata.FileMetaData) *Kademlia {
	instance := &Kademlia{nil, RT, MBL, MD}
	storagePath, _ = filepath.Abs("../storage/")
	downLoadPath, _ = filepath.Abs("../downloads/")

	fmt.Println("Starting timed jobs")
	instance.scheduleMessageBufferListGarbageCollect()
	fmt.Println("RPC timeout garbage collection started!")
	instance.scheduleIdleBucketReExploration()
	fmt.Println("Idle bucket re-exploration started!")
	instance.scheduleFileRepublish()
	fmt.Println("File republishing started!")
	instance.scheduleCacheExpiredFileDeletion()
	fmt.Println("Cache expired file deletion started!")
	fmt.Println("All timed jobs started successfully!")
	return instance
}

func (kademlia *Kademlia) SetNetworkHandler(handler NetworkHandler) {
	kademlia.network = handler
}

func (kademlia *Kademlia) lookupProcedure(procedureType int, target *d7024e.KademliaID) ([]d7024e.Contact, *d7024e.Contact, bool) {

	rTable := kademlia.routingTable
	candids := &d7024e.ContactCandidates{}
	queriedAddresses := make(map[string]bool)

	//select own k known closest node for first iteration
	candids.Append(rTable.FindClosestContacts(target, valueK))
	candids.Sort()

	//Variables used in find value
	var fileHosts []*d7024e.Contact
	fileWasFound := false

	var chans []*chan kademliaMessage

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
						chans = append(chans, &ch)
						go kademlia.lookupSubProcedure(contact, target, procedureType, ch)
						queried = queried + 1
						queriedAddresses[contact.Address] = true
						fmt.Println("started query among alpha to number ", queried)
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
				fmt.Println("started new waittime")
				timeWaited := time.Since(startTime)
				allowedTimeouts := int(timeWaited.Seconds()) / incrementalLimit

				for i := startIndex; i < len(chans); i++ {
					select {
					case x, ok := <-*chans[i]:
						if ok {
							if x.returnType == returnContacts {
								fmt.Println("got back contacts")
								for i, _ := range x.contacts {
									x.contacts[i].CalcDistance(target)
								}
								fmt.Println("appending ", len(x.contacts))
								candids.Append(x.contacts)
							} else if x.returnType == returnHasValue {
								x.contacts[0].CalcDistance(target)
								fileHosts = append(fileHosts, &x.contacts[0])
								if !fileWasFound {
									fileWasFound = true
								}
							} else if x.returnType == returnTimedOut {
								candids.RemoveContact(x.contacts[0].ID)
							}
						}
					default:
						allowedTimeouts--
					}
				}

				if fileWasFound {
					candids.RemoveContact(rTable.Me.ID)
					for _, c := range fileHosts {
						candids.RemoveContact(c.ID)
					}
					candids.Sort()
					distinctContacts := candids.GetDistinctContacts(valueK)
					return distinctContacts, fileHosts[0], fileWasFound
				}

				if allowedTimeouts < 0 {
					time.Sleep(1000 * time.Millisecond)
				} else {
					break
				}
			}

			//go trhough all channels, save responses and remove clsoed channels
			fmt.Println("going through all channels")
			for i := len(chans) - 1; i >= 0; i-- {
				select {
				case x, ok := <-*chans[i]:
					if ok {
						if x.returnType == returnContacts {
							for i, _ := range x.contacts {
								x.contacts[i].CalcDistance(target)
							}
							candids.Append(x.contacts)
						} else if x.returnType == returnHasValue {
							x.contacts[0].CalcDistance(target)
							fileHosts = append(fileHosts, &x.contacts[0])
							if !fileWasFound {
								fileWasFound = true
							}
						} else if x.returnType == returnTimedOut {
							candids.RemoveContact(x.contacts[0].ID)
						}
					}
					chans = append(chans[:i], chans[i+1:]...)
				default:
					continue
				}
			}

			if fileWasFound {
				candids.RemoveContact(rTable.Me.ID)
				for _, c := range fileHosts {
					candids.RemoveContact(c.ID)
				}
				candids.Sort()
				distinctContacts := candids.GetDistinctContacts(valueK)
				return distinctContacts, fileHosts[0], fileWasFound
			}

			//save k closest distinct contacts for next iteration
			candids.RemoveContact(rTable.Me.ID)
			candids.Sort()
			distinctContacts := candids.GetDistinctContacts(valueK)
			candids := &d7024e.ContactCandidates{}
			candids.Append(distinctContacts)

		} else { // last iteration gave no closer nodes, query all
			fmt.Println("quering all")
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
					chans = append(chans, &ch)
					fmt.Println("quering ", contact.ID.String())
					go kademlia.lookupSubProcedure(contact, target, procedureType, ch)
				}
			}

			//wait for responses or time-outs for all open queries
			fmt.Println("waiting for response from all")
			for {
				if len(chans) == 0 {
					break
				}

				for i := len(chans) - 1; i >= 0; i-- {
					select {
					case x, ok := <-*chans[i]:
						if ok {
							if x.returnType == returnContacts {
								fmt.Println("got response, ", i, " remaining")
								for i, _ := range x.contacts {
									x.contacts[i].CalcDistance(target)
								}
								candids.Append(x.contacts)
							} else if x.returnType == returnHasValue {
								x.contacts[0].CalcDistance(target)
								fileHosts = append(fileHosts, &x.contacts[0])
								if !fileWasFound {
									fileWasFound = true
								}
							} else if x.returnType == returnTimedOut {
								candids.RemoveContact(x.contacts[0].ID)
							}
						}
						chans = append(chans[:i], chans[i+1:]...)
					default:
						continue
					}
				}

				if fileWasFound {
					candids.RemoveContact(rTable.Me.ID)
					for _, c := range fileHosts {
						candids.RemoveContact(c.ID)
					}
					candids.Sort()
					distinctContacts := candids.GetDistinctContacts(valueK)
					return distinctContacts, fileHosts[0], fileWasFound
				}

				time.Sleep(1000 * time.Millisecond)
			}
		}

	}
	fmt.Println("returning")
	candids.RemoveContact(rTable.Me.ID)
	candids.Sort()
	distinctContacts := candids.GetDistinctContacts(valueK)
	return distinctContacts, nil, fileWasFound

}

func (kademlia *Kademlia) lookupSubProcedure(target d7024e.Contact, toFind *d7024e.KademliaID, lookupType int, ch chan<- kademliaMessage) {
	//create and add its own messageBuffer to the singleton list
	rpcID := d7024e.NewRandomKademliaID()
	mBuffer := messageBufferList.NewMessageBuffer(rpcID)
	mBufferList := kademlia.MBList
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
		var contacts rpc.ClosestNodes
		json.Unmarshal(message.RpcData, &contacts)
		retMessage := kademliaMessage{returnContacts, contacts.Closest}
		fmt.Println("returning contacts ", len(contacts.Closest))
		ch <- retMessage
		close(ch)
	} else if message.RpcType == rpc.HAS_VALUE {
		//return target so main routine knows which contact has the file
		contacts := []d7024e.Contact{target}
		retMessage := kademliaMessage{returnHasValue, contacts}
		fmt.Println("returning hasfile")
		ch <- retMessage
		close(ch)
	} else {
		//return target so main routine knows which contact has timed out
		contacts := []d7024e.Contact{target}
		retMessage := kademliaMessage{returnTimedOut, contacts}
		fmt.Println("returning timed out")
		ch <- retMessage
		close(ch)
	}
}

func (kademlia *Kademlia) LookupContact(target *d7024e.KademliaID) (closest []d7024e.Contact) {
	closest, _, _ = kademlia.lookupProcedure(procedureContacts, target)
	return
}

func (kademlia *Kademlia) LookupData(id string) (filePath string, closest []d7024e.Contact, fileWasFound bool) {
	fileHash := d7024e.NewKademliaID(id)
	closest, fileHost, fileWasFound := kademlia.lookupProcedure(procedureValue, fileHash)
	if fileWasFound {
		url := fileHost.Address + "/storage/" + id
		filePath = downLoadPath + "/" + id
		kademlia.network.FetchFile(url, filePath)
		if len(closest) > 0 {
			kademlia.sendStoreMessage(&closest[0], d7024e.NewRandomKademliaID(), fileHash, fileHost.Address)
		}
	}
	return
}

func (kademlia *Kademlia) StoreFile(filePath string) {
	file, err := os.Open(filePath)
	fmt.Println()
	defer file.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("opened file ", filePath)

	h := sha1.New()
	_, err = io.Copy(h, file)
	if err != nil {
		fmt.Println(err)
		return
	}
	file.Seek(0, 0)

	fmt.Println("hashed file")

	hash := h.Sum(nil)
	newFileName := hex.EncodeToString(hash)
	newPath := storagePath + "/" + newFileName
	destination, deserr := os.Create(newPath)
	defer destination.Close()
	if deserr != nil {
		fmt.Println(deserr)
		return
	}

	_, wrierr := io.Copy(destination, file)
	io.Copy(os.Stdout, file)
	//fmt.Println(bytes, " bytes copied")
	if wrierr != nil {
		fmt.Println(wrierr)
		return
	}

	fmt.Println("copied file")

	kademliaHash := d7024e.NewKademliaID(newFileName)
	fmt.Println("metadata")
	kademlia.MetaData.AddFile(newPath, newFileName, true, kademlia.calcTimeToLive(kademliaHash))

	fmt.Println("getting closest nodes")
	closest, _, _ := kademlia.lookupProcedure(procedureContacts, kademliaHash)

	for _, c := range closest {
		fmt.Println("sending store to ", c.ID.String())
		kademlia.sendStoreMessage(&c, d7024e.NewRandomKademliaID(), kademliaHash, rpc.SENDER)
	}
	fmt.Println("Sent store RPC")
}

func (kademlia *Kademlia) Join(ip string, port int) bool {
	rt := kademlia.routingTable
	fmt.Println(fmt.Sprintf("joining %q on port %d", ip, port))
	bootstrapContact := d7024e.NewContact(d7024e.NewRandomKademliaID(), fmt.Sprintf("%s:%d", ip, port))

	retry := 1
	for retry < 4 {
		rpcID := d7024e.NewRandomKademliaID()
		mBuffer := messageBufferList.NewMessageBuffer(rpcID)
		mBufferList := kademlia.MBList
		mBufferList.AddMessageBuffer(mBuffer)
		fmt.Println("Trying to connect. Try number ", retry)
		kademlia.sendFindContactMessage(&bootstrapContact, rt.Me.ID, rpcID)

		//wait until a response is
		message := <-mBuffer.MessageChannel

		if message.RpcType == rpc.CLOSEST_NODES {

			var contacts rpc.ClosestNodes = rpc.ClosestNodes{}
			json.Unmarshal(message.RpcData, &contacts)
			fmt.Println("adding contacts ", len(contacts.Closest))
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

func (kademlia *Kademlia) addContact(contact *d7024e.Contact) {

	rt := kademlia.routingTable
	pingContact, inserted := rt.AddContact(*contact)

	if !inserted {
		go func() {
			iPingContact := pingContact
			iInserted := inserted
			iContact := contact
			for {
				rpcID := d7024e.NewRandomKademliaID()
				mBuffer := messageBufferList.NewMessageBuffer(rpcID)
				mBufferList := kademlia.MBList
				mBufferList.AddMessageBuffer(mBuffer)

				kademlia.sendPingMessage(iPingContact, rpcID)

				fmt.Println("sent ping message from bucket")

				message := <-mBuffer.MessageChannel
				fmt.Println("ping executed")

				if message.RpcType == rpc.PONG {
					fmt.Println("ping responded successfully, node alive.")
					rt.AddContact(*iPingContact)
					return
				} else {
					fmt.Println("ping without response, node dead.")
					iPingContact, iInserted = rt.ReplaceLastSeenNode(*iPingContact, *iContact)
					if iInserted {

						return
					}
				}
			}
		}()
	}

}

func (kademlia *Kademlia) calcTimeToLive(fileId *d7024e.KademliaID) (ttl time.Duration) {
	exponent := float64(kademlia.routingTable.GetBucketIndex(fileId))
	ttl = time.Duration(coefficentTTL * math.Exp(exponent))
	return
}
