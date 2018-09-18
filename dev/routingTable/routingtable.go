package routingTable

import (
	"sync"

	"Kademlia_D7024E/dev/d7024e"
)

var instance *routingTable
var once sync.Once

// RoutingTable definition
// keeps a refrence contact of me and an array of buckets
type routingTable struct {
	me      d7024e.Contact
	buckets [d7024e.IDLength * 8]*d7024e.Bucket
}

func GetInstance() *routingTable {
	once.Do(func() {
		instance = newRoutingTable()
	})
	return instance
}

// NewRoutingTable returns a new instance of a RoutingTable
func newRoutingTable() *routingTable {
	routingTable := &routingTable{}
	for i := 0; i < d7024e.IDLength*8; i++ {
		routingTable.buckets[i] = d7024e.NewBucket()
	}
	return routingTable
}

// AddContact add a new contact to the correct Bucket
func (routingTable *routingTable) AddContact(contact d7024e.Contact) {
	bucketIndex := routingTable.getBucketIndex(contact.ID)
	bucket := routingTable.buckets[bucketIndex]
	bucket.AddContact(contact)
}

// FindClosestContacts finds the count closest Contacts to the target in the RoutingTable
func (routingTable *routingTable) FindClosestContacts(target *d7024e.KademliaID, count int) []d7024e.Contact {
	var candidates d7024e.ContactCandidates
	bucketIndex := routingTable.getBucketIndex(target)
	bucket := routingTable.buckets[bucketIndex]

	candidates.Append(bucket.GetContactAndCalcDistance(target))

	for i := 1; (bucketIndex-i >= 0 || bucketIndex+i < d7024e.IDLength*8) && candidates.Len() < count; i++ {
		if bucketIndex-i >= 0 {
			bucket = routingTable.buckets[bucketIndex-i]
			candidates.Append(bucket.GetContactAndCalcDistance(target))
		}
		if bucketIndex+i < d7024e.IDLength*8 {
			bucket = routingTable.buckets[bucketIndex+i]
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
func (routingTable *routingTable) getBucketIndex(id *d7024e.KademliaID) int {
	distance := id.CalcDistance(routingTable.me.ID)
	for i := 0; i < d7024e.IDLength; i++ {
		for j := 0; j < 8; j++ {
			if (distance[i]>>uint8(7-j))&0x1 != 0 {
				return i*8 + j
			}
		}
	}

	return d7024e.IDLength*8 - 1
}
