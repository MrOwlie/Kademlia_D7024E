package kademlia

import (
	"sync"

	"Kademlia_D7024E/dev/d7024e"
	"Kademlia_D7024E/dev/network"
	"Kademlia_D7024E/dev/routingTable"
	"time"
)

type kademlia struct {
}

type kademliaMessage struct {
	returnType int
	contacts   []d7024e.Contact
}

var instance *kademlia
var once sync.Once

const alpha int = 3
const valueK int = 20

const procedureContacts = 1
const procedureValue = 2
const returnContacts = 3
const returnHasValue = 4

func GetInstance() *kademlia {
	once.Do(func() {
		instance = &kademlia{}
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
						go kademlia.lookupSubProcedure(procedureType, ch)
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
								close(chans[i])
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
							close(chans[i])
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
					go kademlia.lookupSubProcedure(procedureType, ch)
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
								//candids.Append(x.contacts)
								close(chans[i])
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

func (kademlia *kademlia) lookupSubProcedure(lookupType int, ch chan<- kademliaMessage) {
	net := network.GetInstance()
	rpcId := d7024e.NewRandomKademliaID()
}

func (kademlia *kademlia) requestFile(file *d7024e.KademliaID, target d7024e.Contact) {

}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func (kademlia *kademlia) LookupData(hash string) {
	// TODO
}

func (kademlia *kademlia) Store(data []byte) {
	// TODO
}

func (kademlia *kademlia) Join(ip string, port int) {
	// TODO
}

func (kademlia *kademlia) ReturnLookupContact(target *d7024e.Contact) {
	// TODO
}

func (kademlia *kademlia) ReturnLookupData(hash string) {
	// TODO
}
