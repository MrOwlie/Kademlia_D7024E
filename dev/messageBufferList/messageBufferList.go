package messageBufferList

import (
	"sync"
	"../d7024e"
	"container/list"
)

type messageBufferList struct {
	list []*messageBuffer
}

var instance *messageBufferList
var once sync.Once

func GetInstance() *messageBufferList {
    once.Do(func() {
		instance = &messageBufferList{}
		instance.list = list.New()
    })
    return instance
}

func (mbList *messageBufferList) AddMessageBuffer(mb *messageBuffer){
	append(mbList.list, mb)
}

func (mbList *messageBufferList) GetMessageBuffer(id *KademliaID) *messageBuffer, bool{
	for _, element := range mbList.list{
		if element.RPCID.Equals(id) {
			return element, true
		}
	}
	return nil, false
}

func (mbList *messageBufferList) DeleteMessageBuffer(id *KademliaID) bool{
	for i, element := range mbList.list{
		if element.RPCID.Equals(id) {
			copy(mbList.list[i:], mbList.list[i+1:])
			mbList.list[len(mbList.list)-1] = nil // or the zero value of T
			mbList.list = mbList.list[:len(mbList.list)-1]

			return true
		}
	}
	return false
}