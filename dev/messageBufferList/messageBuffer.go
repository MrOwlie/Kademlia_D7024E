package messageBufferList

import (
	"Kademlia_D7024E/dev/d7024e"
	"Kademlia_D7024E/dev/rpc"
	"sync"
	"time"
)

const messageBufferTimeout = 5

type messageBuffer struct {
	RPCID          *d7024e.KademliaID
	waitGroup      sync.WaitGroup
	message        rpc.Message
	latestResponse time.Time
}

func newMessageBuffer(id *d7024e.KademliaID) *messageBuffer {
	mb := &messageBuffer{}
	mb.RPCID = id
	return mb
}

func (mb *messageBuffer) AppendMessage(message rpc.Message) {
	mb.message = message
	mb.latestResponse = time.Now()
	mb.waitGroup.Done()
}

func (mb *messageBuffer) extractMessage() (r_message rpc.Message) {
	r_message = mb.message
	return
}

func (mb *messageBuffer) hasExpired() bool {
	elapsed := time.Since(mb.latestResponse)
	if elapsed.Minutes() > messageBufferTimeout {
		return true
	}
	return false
}
