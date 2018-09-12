package kademlia

import (
	"sync"
	"../d7024e"
)

type kademlia struct {
}

var instance *kademlia
var once sync.Once

func GetInstance() *kademlia {
    once.Do(func() {
        instance = &kademlia{}
    })
    return instance
}

func (kademlia *kademlia) LookupContact(target *Contact) {
	// TODO
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

func (kademlia *kademlia) ReturnLookupContact(target *Contact) {
	// TODO
}

func (kademlia *kademlia) ReturnLookupData(hash string) {
	// TODO
}
