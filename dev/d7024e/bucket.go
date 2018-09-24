package d7024e

import (
	"container/list"
	"fmt"
	"sync"
	"time"

	"../d7024e"
	"../kademlia"
	"../messageBufferList"
)

const bucketSize = 20
const timeForRefresh = 60

// bucket definition
// contains a List
type Bucket struct {
	list         *list.List
	mutex        sync.Mutex
	latestLookup time.Time
}

// newBucket returns a new instance of a bucket
func NewBucket() *Bucket {
	bucket := &Bucket{}
	bucket.list = list.New()
	return bucket
}

// AddContact adds the Contact to the front of the bucket
// or moves it to the front of the bucket if it already existed
func (bucket *Bucket) AddContact(contact Contact) {
	var element *list.Element
	bucket.mutex.Lock()
	defer bucket.mutex.Unlock()
	bucket.latestLookup = time.Now()

	for e := bucket.list.Front(); e != nil; e = e.Next() {
		nodeID := e.Value.(Contact).ID

		if (contact).ID.Equals(nodeID) {
			element = e
		}
	}

	if element == nil {
		if bucket.list.Len() < bucketSize {
			bucket.list.PushBack(contact)
		} else {
			kademlia := kademlia.GetInstance()
			rpcID := d7024e.NewRandomKademliaID()
			pingContact := bucket.list.Front()
			mBuffer := messageBufferList.NewMessageBuffer(rpcID)
			mBufferList := messageBufferList.GetInstance()
			mBufferList.AddMessageBuffer(mBuffer)

			kademlia.sendPingMessage(&pingContact, rpcID)

			fmt.Println("sent ping message from bucket")

			mBuffer.WaitForResponse()
			message := mBuffer.ExtractMessage()
			fmt.Println("ping executed")

			if (message != nil) || (message.rpcType == rpc.PONG) {
				fmt.Println("ping responded successfully, node alive.")
				bucket.list.MoveToBack(&pingContact)
			} else {
				fmt.Println("ping without response, node dead.")
				bucket.list.Remove(&pingContact)
				bucket.list.PushBack(contact)
			}

		}
	} else {
		bucket.list.MoveToFront(element)
	}
}

// GetContactAndCalcDistance returns an array of Contacts where
// the distance has already been calculated
func (bucket *Bucket) GetContactAndCalcDistance(target *KademliaID) []Contact {
	var contacts []Contact
	bucket.mutex.Lock()
	defer bucket.mutex.Unlock()
	bucket.latestLookup = time.Now()

	for elt := bucket.list.Front(); elt != nil; elt = elt.Next() {
		contact := elt.Value.(Contact)
		contact.CalcDistance(target)
		contacts = append(contacts, contact)
	}

	return contacts
}

// Len return the size of the bucket
func (bucket *Bucket) Len() int {
	return bucket.list.Len()
}

func (bucket *Bucket) NeedsRefresh() bool {
	elapsed := time.Since(bucket.latestLookup)
	return (elapsed.Minutes() > timeForRefresh)
}
