package messageBufferList

import (
	"sync"
	"../d7024e"
	"time"
)

const messageBufferTimeout = 5

type messageBuffer struct {
	RPCID *KademliaID
	waitGroup sync.WaitGroup
	messageList []string
	latestResponse time.Time
}

func newMessageBuffer(id *KademliaID) *messageBuffer{
	mb := &messageBuffer{}
	mb.RPCID = id
}

func (mb *messageBuffer) appendMessages(messages []string){
	for _, element := range messages{
		append(mb.messageList, element)
	}
	mb.latestResponse = time.Now()
	mb.waitGroup.Done()
}

func (mb *messageBuffer) extractMessages() []string{
	temp := []string
	copy(temp, mb.messageList)
	mb.messageList = nil
}

func (mb *messageBuffer) hasExpired() bool{
	elapsed := time.Since(mb.latestResponse)
	if time.Minutes(elapsed) > messageBufferTimeout {
		return true
	}
	return false
}
