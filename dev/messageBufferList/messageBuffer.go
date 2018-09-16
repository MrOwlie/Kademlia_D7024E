package messageBufferList

import (
	"../d7024e"
	"../network"
	"sync"
	"time"
)

const messageBufferTimeout = 5

type messageBuffer struct {
	RPCID          *d7024e.KademliaID
	waitGroup      sync.WaitGroup
	message        network.Message
	latestResponse time.Time
}

func newMessageBuffer(id *d7024e.KademliaID) *messageBuffer {
	mb := &messageBuffer{}
	mb.RPCID = id
	return mb
}

func (mb *messageBuffer) appendMessage(message network.Message) {
	mb.message = message
	mb.latestResponse = time.Now()
	mb.waitGroup.Done()
}

func (mb *messageBuffer) extractMessages() r_message network.Message {
	r_message = messageBuffer.message
	return 
}

func (mb *messageBuffer) hasExpired() bool {
	elapsed := time.Since(mb.latestResponse)
	if elapsed.Minutes() > messageBufferTimeout {
		return true
	}
	return false
}
