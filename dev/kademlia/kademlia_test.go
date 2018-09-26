package kademlia

import (
	"encoding/json"
	"fmt"
	"sync"
	"testing"

	"../d7024e"
	"../routingTable"
	"../rpc"
)

var doNotCareID *d7024e.KademliaID = d7024e.NewKademliaID("F000000000000000000000000000000000000000")

type testNetworkControl struct {
	SentMsg     rpc.Message
	ReturnMsg   rpc.Message
	FromAddress string
}

type testNetwork struct {
	CurrentTest *testing.T
	CheckList   []testNetworkControl
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
	sentMessageControl := checkData.SentMsg
	returnMessage := checkData.ReturnMsg

	assertEqual(net.CurrentTest, sentMessage.RpcType, sentMessageControl.RpcType)
	if !sentMessageControl.SenderId.Equals(doNotCareID) {
		fmt.Println("testing sender ID")
		assertEqual(net.CurrentTest, sentMessageControl.SenderId.Equals(&sentMessage.SenderId), true)
	}

	switch sentMessage.RpcType {
	case rpc.FIND_NODE:
		var sentPayload rpc.FindNode
		var controlPayload rpc.FindNode
		json.Unmarshal(sentMessage.RpcData, &sentPayload)
		json.Unmarshal(sentMessageControl.RpcData, &controlPayload)

		fmt.Println("testing FIND NODE ID: should be ", controlPayload.NodeId.String(), ". is ", sentPayload.NodeId.String())
		assertEqual(net.CurrentTest, sentPayload.NodeId.Equals(&controlPayload.NodeId), true)
		fmt.Println("testing FIND_NODE successful!")
	}

	returnMessage.RpcId = sentMessage.RpcId
	d, _ := json.Marshal(returnMessage)

	GetInstance().HandleIncomingRPC(d, checkData.FromAddress)
}



func TestFindNode(t *testing.T) {
	target := d7024e.NewContact(d7024e.NewKademliaID("FFFFFFFFF0000000000000000000000000000000"), "localhost:8000")
	expectedResponse := []d7024e.Contacts{					
		d7024e.NewContact(d7024e.NewKademliaID("FFFFFFFFF0000000000000000000000000000001"), "localhost:8001"),
		d7024e.NewContact(d7024e.NewKademliaID("FFFFFFFFF0000000000000000000000000000002"), "localhost:8002"),
		d7024e.NewContact(d7024e.NewKademliaID("FFFFFFFFF0000000000000000000000000000003"), "localhost:8003"),
		d7024e.NewContact(d7024e.NewKademliaID("FFFFFFFFF0000000000000000000000000000004"), "localhost:8004"),
		d7024e.NewContact(d7024e.NewKademliaID("FFFFFFFFF0000000000000000000000000000005"), "localhost:8005"),
		d7024e.NewContact(d7024e.NewKademliaID("FFFFFFFFF0000000000000000000000000000006"), "localhost:8006"),
		d7024e.NewContact(d7024e.NewKademliaID("FFFFFFFFF0000000000000000000000000000007"), "localhost:8007"),
		d7024e.NewContact(d7024e.NewKademliaID("FFFFFFFFF0000000000000000000000000000008"), "localhost:8008"),
		d7024e.NewContact(d7024e.NewKademliaID("FFFFFFFFF0000000000000000000000000000009"), "localhost:8009"),
		d7024e.NewContact(d7024e.NewKademliaID("FFFFFFFFF000000000000000000000000000000A"), "localhost:8010"),
		d7024e.NewContact(d7024e.NewKademliaID("FFFFFFFFF000000000000000000000000000000B"), "localhost:8011"),
		d7024e.NewContact(d7024e.NewKademliaID("FFFFFFFFF000000000000000000000000000000C"), "localhost:8012"),
		d7024e.NewContact(d7024e.NewKademliaID("FFFFFFFFF000000000000000000000000000000D"), "localhost:8013"),
		d7024e.NewContact(d7024e.NewKademliaID("FFFFFFFFF000000000000000000000000000000E"), "localhost:8014"),
		d7024e.NewContact(d7024e.NewKademliaID("FFFFFFFFF000000000000000000000000000000F"), "localhost:8015"),
		d7024e.NewContact(d7024e.NewKademliaID("FFFFFFFFF0000000000000000000000000000010"), "localhost:8016"),
		d7024e.NewContact(d7024e.NewKademliaID("FFFFFFFFF0000000000000000000000000000011"), "localhost:8017"),
		d7024e.NewContact(d7024e.NewKademliaID("FFFFFFFFF0000000000000000000000000000012"), "localhost:8018"),
		d7024e.NewContact(d7024e.NewKademliaID("FFFFFFFFF0000000000000000000000000000013"), "localhost:8019"),
		d7024e.NewContact(d7024e.NewKademliaID("FFFFFFFFF0000000000000000000000000000014"), "localhost:8020"),
	}


}

func TestJoin(t *testing.T) {
	kadem := GetInstance()
	rt := routingTable.GetInstance()
	cList := []testNetworkControl{
		testNetworkControl{
			rpc.Message{rpc.FIND_NODE, *doNotCareID, *rt.Me.ID, helperReturnMarshal(rpc.FindNode{*rt.Me.ID})},
			rpc.Message{rpc.CLOSEST_NODES, *doNotCareID, *d7024e.NewKademliaID("FFFFFFFF00000000000000000000000000000000"), helperReturnMarshal(rpc.ClosestNodes{
				[]d7024e.Contact{
					d7024e.NewContact(d7024e.NewKademliaID("FFFFFFFFF0000000000000000000000000000001"), "localhost:8000"),
					d7024e.NewContact(d7024e.NewKademliaID("FFFFFFFFF0000000000000000000000000000002"), "localhost:8001"),
					d7024e.NewContact(d7024e.NewKademliaID("FFFFFFFFF0000000000000000000000000000003"), "localhost:8002"),
					d7024e.NewContact(d7024e.NewKademliaID("FFFFFFFFF0000000000000000000000000000004"), "localhost:8003"),
					d7024e.NewContact(d7024e.NewKademliaID("FFFFFFFFF0000000000000000000000000000005"), "localhost:8004"),
					d7024e.NewContact(d7024e.NewKademliaID("FFFFFFFFF0000000000000000000000000000006"), "localhost:8005"),
				},
			}),
			},
			"10.10.10.12:1000",
		},
	}
	net := testNetwork{t, nil}
	net.CheckList = cList
	kadem.SetNetworkHandler(&net)
	var wg sync.WaitGroup
	wg.Add(1)
	kadem.Join("1.1.1.1", 1000, &wg)
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
