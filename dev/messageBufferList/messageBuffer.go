package messageBufferList

import (
	"sync"
	"time"

	"../d7024e"
)

const messageBufferTimeout = 5

type messageBuffer struct {
	RPCID          *d7024e.KademliaID
	waitGroup      sync.WaitGroup
	messageList    []string
	latestResponse time.Time
}

func newMessageBuffer(id *d7024e.KademliaID) *messageBuffer {
	mb := &messageBuffer{}
	mb.RPCID = id
	return mb
}

func (mb *messageBuffer) appendMessages(messages []string) {
	for _, element := range messages {
		mb.messageList = append(mb.messageList, element)
	}
	mb.latestResponse = time.Now()
	mb.waitGroup.Done()
}

func (mb *messageBuffer) extractMessages() []string {
	var temp []string
	copy(temp, mb.messageList)
	mb.messageList = nil
	return temp
}

func (mb *messageBuffer) hasExpired() bool {
	elapsed := time.Since(mb.latestResponse)
	if elapsed.Minutes() > messageBufferTimeout {
		return true
	}
	return false
}
