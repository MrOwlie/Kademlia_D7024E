package routingTable

import (
	"fmt"
	"sync"

	"../d7024e"
)

var instance *routingTable
var once sync.Once

// RoutingTable definition
// keeps a refrence contact of me and an array of buckets
type routingTable struct {
	Me      d7024e.Contact
	Buckets [d7024e.IDLength * 8]*d7024e.Bucket
}

func GetInstance() *routingTable {
	once.Do(func() {
		instance = newRoutingTable()
		instance.Me = d7024e.Contact{}
		instance.Me.ID = d7024e.NewRandomKademliaID()
		fmt.Println("my id is: ", instance.Me.ID)
	})
	return instance
}

// NewRoutingTable returns a new instance of a RoutingTable
func newRoutingTable() *routingTable {
	routingTable := &routingTable{}
	for i := 0; i < d7024e.IDLength*8; i++ {
		routingTable.Buckets[i] = d7024e.NewBucket()
	}
	return routingTable
}

// AddContact add a new contact to the correct Bucket
func (routingTable *routingTable) AddContact(contact d7024e.Contact) {
	bucketIndex := routingTable.GetBucketIndex(contact.ID)
	bucket := routingTable.Buckets[bucketIndex]
	if !routingTable.Me.ID.Equals(contact.ID) {
		bucket.AddContact(contact)
	}
}

// FindClosestContacts finds the count closest Contacts to the target in the RoutingTable
func (routingTable *routingTable) FindClosestContacts(target *d7024e.KademliaID, count int) []d7024e.Contact {
	var candidates d7024e.ContactCandidates
	bucketIndex := routingTable.GetBucketIndex(target)
	bucket := routingTable.Buckets[bucketIndex]

	candidates.Append(bucket.GetContactAndCalcDistance(target))

	for i := 1; (bucketIndex-i >= 0 || bucketIndex+i < d7024e.IDLength*8) && candidates.Len() < count; i++ {
		if bucketIndex-i >= 0 {
			bucket = routingTable.Buckets[bucketIndex-i]
			candidates.Append(bucket.GetContactAndCalcDistance(target))
		}
		if bucketIndex+i < d7024e.IDLength*8 {
			bucket = routingTable.Buckets[bucketIndex+i]
			candidates.Append(bucket.GetContactAndCalcDistance(target))
		}
	}

	candidates.Sort()

	if count > candidates.Len() {
		count = candidates.Len()
	}

	return candidates.GetContacts(count)
}

// getBucketIndex get the correct Bucket index for the KademliaID
func (routingTable *routingTable) GetBucketIndex(id *d7024e.KademliaID) int {
	distance := id.CalcDistance(routingTable.Me.ID)
	for i := 0; i < d7024e.IDLength; i++ {
		for j := 0; j < 8; j++ {
			if (distance[i]>>uint8(7-j))&0x1 != 0 {
				return i*8 + j
			}
		}
	}

	return d7024e.IDLength*8 - 1
}

//Returns an array of random kademliaIds, where each kademliaID is in the range of a bucket that needs to be refreshed.
func (routingTable *routingTable) GetRefreshIDs() []*d7024e.KademliaID {

	var idList []*d7024e.KademliaID

	for i := 0; i < len(routingTable.Buckets); i++ {
		if routingTable.Buckets[i].NeedsRefresh() {
			n_fullBytes := i / 8
			var n_bitsInToByte uint64 = uint64(i % 8)
			randomId := d7024e.NewRandomKademliaID() //Used to get a random number of correct size

			for j := 0; j < n_fullBytes; j++ {
				randomId[j] = byte(0)
			}

			mask := byte(255) >> n_bitsInToByte
			randomId[n_fullBytes] &= mask

			//Makes sure that the bit in position 7 - n_bitsInToByte is 1.
			mostSigBit := byte(1) << (7 - n_bitsInToByte)
			randomId[n_fullBytes] |= mostSigBit

			kandemliaId := routingTable.Me.ID.CalcDistance(randomId)
			idList = append(idList, kandemliaId)
		}
	}

	return idList
}
