package kademlia

import (
	"sync"

	"Kademlia_D7024E/dev/d7024e"
	"Kademlia_D7024E/dev/network"
	"Kademlia_D7024E/dev/routingTable"
)

type kademlia struct {
}

var instance *kademlia
var once sync.Once
var selfContact *d7024e.Contact = d7024e.Contact.NewContact(d7024e.NewRandomKademliaID, "localhost");

const alpha int = 3
const valueK int = 20

func GetInstance() *kademlia {
	once.Do(func() {
		instance = &kademlia{}
	})
	return instance
}

func (kademlia *kademlia) LookupContact(target *d7024e.Contact) {
	net := network.GetInstance()
	rTable := routingTable.GetInstance()
	candids := &d7024e.ContactCandidates{}

	candids.Append(rTable.FindClosestContacts(target.ID, valueK))
	candids.Sort()

	//var returnArrays [alpha]*[]d7024e.Contact

	net.SendFindContactMessage(target, target.ID)
}

func (kademlia *kademlia) lookupAlphaContact(target *d7024e.Contact, list *[]*d7024e.Contact, mutex *sync.Mutex) {

}

func (kademlia *kademlia) LookupData(hash string) {
	// TODO
}

func (kademlia *kademlia) Store(data []byte) {
	// TODO
}

func (kademlia *kademlia) Join(ip string, port int) {
	// Lookup self on the ip and port of the known node.
	kademlia.GetInstance.LookupContact(selfContact)
}

func (kademlia *kademlia) ReturnLookupContact(target *d7024e.Contact) {
	// TODO
}

func (kademlia *kademlia) ReturnLookupData(hash string) {
	// TODO
}
