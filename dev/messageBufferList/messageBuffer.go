package messageBufferList

import (
	"time"

	"../d7024e"
	"../rpc"
)

const messageBufferTimeout = 30

type messageBuffer struct {
	RPCID          *d7024e.KademliaID
	MessageChannel chan *rpc.Message
	latestResponse time.Time
}

func NewMessageBuffer(id *d7024e.KademliaID) *messageBuffer {
	mb := &messageBuffer{}
	mb.RPCID = id
	mb.MessageChannel = make(chan *rpc.Message)
	return mb
}

func (mb *messageBuffer) AppendMessage(message *rpc.Message) {
	mb.MessageChannel <- message
	mb.latestResponse = time.Now()
}

//func (mb *messageBuffer) ExtractMessage() rpc.Message {
//	return mb.message
//}

//func (mb *messageBuffer) WaitForResponse() {
//	mb.waitGroup.Add(1)
//	mb.waitGroup.Wait()
//}

func (mb *messageBuffer) hasExpired() bool {
	elapsed := time.Since(mb.latestResponse)
	if elapsed.Seconds() > messageBufferTimeout {
		return true
	}
	return false
}
