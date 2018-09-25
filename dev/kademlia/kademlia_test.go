package kademlia

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"../d7024e"
	"../routingTable"
	"../rpc"
)

var doNotCareID *d7024e.KademliaID = d7024e.NewKademliaID("F000000000000000000000000000000000000000")

type testNetworkControl struct {
	CheckFunction func(rpc.Message)
	//ReturnMsg     rpc.Message
	//FromAddress   string
}

type testNetwork struct {
	CheckList []testNetworkControl
}

func (net *testNetwork) SendMessage(addr string, data *[]byte) {
	checkData := net.CheckList[0]
	if len(net.CheckList) > 1 {
		net.CheckList = net.CheckList[1:]
	} else {
		net.CheckList = nil
	}

	var sentMessage = rpc.Message{}
	json.Unmarshal(*data, &sentMessage)

	checkData.CheckFunction(sentMessage)

}

func (net *testNetwork) FetchFile(a string, b string) error {
	return nil
}

func TestFindNode(t *testing.T) {

}

func TestJoin(t *testing.T) {
	kadem := GetInstance()
	rt := routingTable.GetInstance()
	cList := []testNetworkControl{
		testNetworkControl{
			func(sentMessage rpc.Message) {

				assertEqual(t, sentMessage.RpcType, rpc.FIND_NODE)
				fmt.Println("RPC Type is correct!")
				assertEqual(t, rt.Me.ID.Equals(&sentMessage.SenderId), true)
				fmt.Println("Sender ID is correct!")

				var sentPayload rpc.FindNode
				json.Unmarshal(sentMessage.RpcData, &sentPayload)
				assertEqual(t, sentPayload.NodeId.Equals(rt.Me.ID), true)
				fmt.Println("FIND_NODE ID is correct!")

				returnMessage := rpc.Message{rpc.CLOSEST_NODES, sentMessage.RpcId, *d7024e.NewKademliaID("FFFFFFFFF0000000000000000000000000000000"), helperReturnMarshal(rpc.ClosestNodes{
					[]d7024e.Contact{
						d7024e.NewContact(d7024e.NewKademliaID("FFFFFFFFF0000000000000000000000000000001"), "localhost:8000"),
						d7024e.NewContact(d7024e.NewKademliaID("FFFFFFFFF0000000000000000000000000000002"), "localhost:8001"),
						d7024e.NewContact(d7024e.NewKademliaID("FFFFFFFFF0000000000000000000000000000003"), "localhost:8002"),
						d7024e.NewContact(d7024e.NewKademliaID("FFFFFFFFF0000000000000000000000000000004"), "localhost:8003"),
						d7024e.NewContact(d7024e.NewKademliaID("FFFFFFFFF0000000000000000000000000000005"), "localhost:8004"),
						d7024e.NewContact(d7024e.NewKademliaID("FFFFFFFFF0000000000000000000000000000006"), "localhost:8005"),
					},
				}),
				}
				d, _ := json.Marshal(returnMessage)

				go GetInstance().HandleIncomingRPC(d, "10.10.10.10:1000")

			},
		},
	}
	net := testNetwork{}
	net.CheckList = cList
	kadem.SetNetworkHandler(&net)
	assertEqual(t, kadem.Join("1.1.1.1", 1000), true)
	fmt.Println("Join returns correct value!")

	checkMap := make(map[string]bool)
	checkMap["FFFFFFFFF0000000000000000000000000000000"] = true
	checkMap["FFFFFFFFF0000000000000000000000000000001"] = true
	checkMap["FFFFFFFFF0000000000000000000000000000002"] = true
	checkMap["FFFFFFFFF0000000000000000000000000000003"] = true
	checkMap["FFFFFFFFF0000000000000000000000000000004"] = true
	checkMap["FFFFFFFFF0000000000000000000000000000005"] = true
	checkMap["FFFFFFFFF0000000000000000000000000000006"] = true
	checkMap["localhost:8000"] = true
	checkMap["localhost:8001"] = true
	checkMap["localhost:8002"] = true
	checkMap["localhost:8003"] = true
	checkMap["localhost:8004"] = true
	checkMap["localhost:8005"] = true
	checkMap["10.10.10.10:1000"] = true

	contacts := rt.FindClosestContacts(d7024e.NewKademliaID("FFFFFFFFF0000000000000000000000000000000"), 10)
	assertEqual(t, len(contacts), 7)
	fmt.Println("Correct number of contacts where inserted")
	for _, c := range contacts {
		assertEqual(t, checkMap[strings.ToUpper(c.ID.String())], true)
		assertEqual(t, checkMap[c.Address], true)
		checkMap[c.ID.String()] = false
		checkMap[c.Address] = false
	}
	fmt.Println("Inserted contacts had right values")

}

func helperReturnMarshal(data interface{}) []byte {
	da, _ := json.Marshal(data)
	return da
}

func assertEqual(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Fatalf("%s != %s", a, b)
	}
}
