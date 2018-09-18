package messageBufferList

import (
	"sync"

	"Kademlia_D7024E/dev/d7024e"
)

type messageBufferList struct {
	list []*messageBuffer
}

var instance *messageBufferList
var once sync.Once

func GetInstance() *messageBufferList {
	once.Do(func() {
		instance = &messageBufferList{}
	})
	return instance
}

func (mbList *messageBufferList) AddMessageBuffer(mb *messageBuffer) {
	mbList.list = append(mbList.list, mb)
}

func (mbList *messageBufferList) GetMessageBuffer(id *d7024e.KademliaID) (*messageBuffer, bool) {
	for _, element := range mbList.list {
		if element.RPCID.Equals(id) {
			return element, true
		}
	}
	return nil, false
}

func (mbList *messageBufferList) DeleteMessageBuffer(id *d7024e.KademliaID) bool {
	for i, element := range mbList.list {
		if element.RPCID.Equals(id) {
			copy(mbList.list[i:], mbList.list[i+1:])
			mbList.list[len(mbList.list)-1] = nil // or the zero value of T
			mbList.list = mbList.list[:len(mbList.list)-1]

			return true
		}
	}
	return false
}
