package messageBufferList

import (
	"sync"

	"../d7024e"
	"../rpc"
)

type messageBufferList struct {
	list  []*messageBuffer
	mutex sync.Mutex
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
	mbList.mutex.Lock()
	mbList.list = append(mbList.list, mb)
	mbList.mutex.Unlock()
}

func (mbList *messageBufferList) GetMessageBuffer(id *d7024e.KademliaID) (*messageBuffer, bool) {
	mbList.mutex.Lock()
	for _, element := range mbList.list {
		if element.RPCID.Equals(id) {
			mbList.mutex.Unlock()
			return element, true
		}
	}
	mbList.mutex.Unlock()
	return nil, false
}

func (mbList *messageBufferList) DeleteMessageBuffer(id *d7024e.KademliaID) bool {
	mbList.mutex.Lock()
	for i, element := range mbList.list {
		if element.RPCID.Equals(id) {
			copy(mbList.list[i:], mbList.list[i+1:])
			mbList.list[len(mbList.list)-1] = nil // or the zero value of T
			mbList.list = mbList.list[:len(mbList.list)-1]

			mbList.mutex.Unlock()
			return true
		}
	}
	mbList.mutex.Unlock()
	return false
}

func (mbList *messageBufferList) GarbageCollect() {
	mbList.mutex.Lock()
	for i := len(mbList.list) - 1; i >= 0; i-- {
		if mbList.list[i].hasExpired() {
			mbList.list[i].MessageChannel <- &rpc.Message{RpcType: rpc.TIME_OUT, RpcId: *mbList.list[i].RPCID, SenderId: *d7024e.NewRandomKademliaID(), RpcData: nil}
			copy(mbList.list[i:], mbList.list[i+1:])
			mbList.list[len(mbList.list)-1] = nil // or the zero value of T
			mbList.list = mbList.list[:len(mbList.list)-1]
		}
	}
	mbList.mutex.Unlock()
}
