package routingTable

import (
	"fmt"
	"testing"

	"../d7024e"
)

func TestRoutingTable(t *testing.T) {
	rt := GetInstance()
	rt.AddContact(d7024e.NewContact(d7024e.NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:8000"))
	rt.AddContact(d7024e.NewContact(d7024e.NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:8001"))
	rt.AddContact(d7024e.NewContact(d7024e.NewKademliaID("1111111100000000000000000000000000000000"), "localhost:8002"))
	rt.AddContact(d7024e.NewContact(d7024e.NewKademliaID("1111111200000000000000000000000000000000"), "localhost:8002"))
	rt.AddContact(d7024e.NewContact(d7024e.NewKademliaID("1111111300000000000000000000000000000000"), "localhost:8002"))
	rt.AddContact(d7024e.NewContact(d7024e.NewKademliaID("1111111400000000000000000000000000000000"), "localhost:8002"))
	rt.AddContact(d7024e.NewContact(d7024e.NewKademliaID("2111111400000000000000000000000000000000"), "localhost:8002"))

	contacts := rt.FindClosestContacts(d7024e.NewKademliaID("2111111400000000000000000000000000000000"), 20)
	for i := range contacts {
		fmt.Println(contacts[i].String())
	}
}
