package messageBufferList

import (
	"sync"
	"time"

	"../d7024e"
	"../rpc"
)

const messageBufferTimeout = 5

type messageBuffer struct {
	RPCID          *d7024e.KademliaID
	waitGroup      sync.WaitGroup
	message        rpc.Message
	latestResponse time.Time
}

func NewMessageBuffer(id *d7024e.KademliaID) *messageBuffer {
	mb := &messageBuffer{}
	mb.RPCID = id
	return mb
}

func (mb *messageBuffer) AppendMessage(message rpc.Message) {
	mb.message = message
	mb.latestResponse = time.Now()
	mb.waitGroup.Done()
}

func (mb *messageBuffer) ExtractMessage() rpc.Message {
	return mb.message
}

func (mb *messageBuffer) WaitForResponse() {
	mb.waitGroup.Add(1)
	mb.waitGroup.Wait()
}

func (mb *messageBuffer) hasExpired() bool {
	elapsed := time.Since(mb.latestResponse)
	if elapsed.Minutes() > messageBufferTimeout {
		return true
	}
	return false
}
