package d7024e

import (
	"container/list"
	"sync"
	"time"
)

const bucketSize = 20
const timeForRefresh = 60

// bucket definition
// contains a List
type Bucket struct {
	List         *list.List
	mutex        sync.Mutex
	latestLookup time.Time
}

// newBucket returns a new instance of a bucket
func NewBucket() *Bucket {
	bucket := &Bucket{}
	bucket.List = list.New()
	return bucket
}

// AddContact adds the Contact to the front of the bucket
// or moves it to the front of the bucket if it already existed
func (bucket *Bucket) AddContact(contact Contact) {
	//fmt.Println("contact added: ", contact)
	var element *list.Element
	bucket.mutex.Lock()
	defer bucket.mutex.Unlock()
	bucket.latestLookup = time.Now()

	for e := bucket.List.Front(); e != nil; e = e.Next() {
		nodeID := e.Value.(Contact).ID

		if (contact).ID.Equals(nodeID) {
			element = e
		}
	}

	if element == nil {
		if bucket.List.Len() < bucketSize {
			bucket.List.PushBack(contact)
		}
	} else {
		bucket.List.MoveToBack(element)
	}
}

// GetContactAndCalcDistance returns an array of Contacts where
// the distance has already been calculated
func (bucket *Bucket) GetContactAndCalcDistance(target *KademliaID) []Contact {
	var contacts []Contact
	bucket.mutex.Lock()
	defer bucket.mutex.Unlock()
	bucket.latestLookup = time.Now()

	for elt := bucket.List.Front(); elt != nil; elt = elt.Next() {
		contact := elt.Value.(Contact)
		contact.CalcDistance(target)
		contacts = append(contacts, contact)
	}

	return contacts
}

// Len return the size of the bucket
func (bucket *Bucket) Len() int {
	return bucket.List.Len()
}

func (bucket *Bucket) NeedsRefresh() bool {
	elapsed := time.Since(bucket.latestLookup)
	return (elapsed.Minutes() > timeForRefresh)
}

func (bucket *Bucket) IsFull() bool {
	if bucket.List.Len() < bucketSize {
		return false
	}

	return true

}
