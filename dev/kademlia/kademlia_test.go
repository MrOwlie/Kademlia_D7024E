package kademlia

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"testing"

	"../d7024e"
	"../messageBufferList"
	"../routingTable"
	"../rpc"
)

var doNotCareID *d7024e.KademliaID = d7024e.NewKademliaID("F000000000000000000000000000000000000000")

type testNetworkControl struct {
	CheckFunction func(rpc.Message, string)
	//ReturnMsg     rpc.Message
	//FromAddress   string
}

type testNetwork struct {
	CheckList []testNetworkControl
	sendMutex sync.Mutex
}

func (net *testNetwork) SendMessage(addr string, data *[]byte) {
	net.sendMutex.Lock()
	defer net.sendMutex.Unlock()

	checkData := net.CheckList[0]
	if len(net.CheckList) > 1 {
		net.CheckList = net.CheckList[1:]
	} else {
		net.CheckList = nil
	}

	var sentMessage = rpc.Message{}
	json.Unmarshal(*data, &sentMessage)

	checkData.CheckFunction(sentMessage, addr)

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
			func(sentMessage rpc.Message, addr string) {

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
	fmt.Println("Correct number of contacts where inserted!")
	for _, c := range contacts {
		assertEqual(t, checkMap[strings.ToUpper(c.ID.String())], true)
		assertEqual(t, checkMap[c.Address], true)
		checkMap[c.ID.String()] = false
		checkMap[c.Address] = false
	}
	fmt.Println("Inserted contacts had right values!")

}

func TestBucketReExploration(t *testing.T) {
	kadem := GetInstance()
	rt := routingTable.GetInstance()
	rt.Me.ID = d7024e.NewKademliaID("0000000000000000000000000000000000000000")
	rt.AddContact(d7024e.NewContact(d7024e.NewKademliaID("0000000000000000000000000000000000000006"), "localhost:8000"))
	checkMap := make(map[int]bool)
	for i := 0; i < 157; i++ {
		checkMap[i] = true
	}
	checkMap[158] = true
	checkMap[159] = true
	var callsMade sync.WaitGroup
	callsMade.Add(159)

	var cList []testNetworkControl
	for a := 0; a < 159; a++ {
		cList = append(cList, testNetworkControl{
			func(sentMessage rpc.Message, addr string) {

				var sentPayload rpc.FindNode
				json.Unmarshal(sentMessage.RpcData, &sentPayload)
				index := rt.GetBucketIndex(&sentPayload.NodeId)

				assertEqual(t, sentMessage.RpcType, rpc.FIND_NODE)
				fmt.Println("RPC Type is correct for index ", index)
				assertEqual(t, rt.Me.ID.Equals(&sentMessage.SenderId), true)
				fmt.Println("Sender ID is correct for index ", index)

				assertEqual(t, checkMap[index], true)
				checkMap[index] = false
				fmt.Println("FIND_NODE ID is within an unqueried bucket, index ", index)

				callsMade.Done()

			},
		})
	}

	net := testNetwork{}
	net.CheckList = cList
	kadem.SetNetworkHandler(&net)
	kadem.IdleBucketReExploration()
	callsMade.Wait()

}

func TestAddContact(t *testing.T) {
	rt := routingTable.GetInstance()
	kadem := GetInstance()
	rt.Me.ID = d7024e.NewKademliaID("FFFF000000000000000000000000000000000000")

	checkdata := make(map[string]bool)

	init := func(id string, ip string) *d7024e.Contact {
		v := d7024e.NewContact(d7024e.NewKademliaID(id), ip)
		checkdata[id] = true
		checkdata[ip] = true
		return &v
	}

	var cList []testNetworkControl
	cList = append(cList, testNetworkControl{
		func(sentMessage rpc.Message, addr string) {

			assertEqual(t, addr, "10.10.10.10:1000")
			fmt.Println("Last seen node is queried!")

			assertEqual(t, sentMessage.RpcType, rpc.PING)
			fmt.Println("RPC Type is correct!")
			assertEqual(t, rt.Me.ID.Equals(&sentMessage.SenderId), true)
			fmt.Println("Sender ID is correct!")

			returnMessage := rpc.Message{rpc.PONG, sentMessage.RpcId, *d7024e.NewKademliaID("F000000000000000000000000000000000000000"), []byte{byte(0)}}
			d, _ := json.Marshal(returnMessage)

			go GetInstance().HandleIncomingRPC(d, addr)

		},
	})
	cList = append(cList, testNetworkControl{
		func(sentMessage rpc.Message, addr string) {

			assertEqual(t, addr, "10.10.10.11:1000")
			fmt.Println("Last seen node is queried!")

			assertEqual(t, sentMessage.RpcType, rpc.PING)
			fmt.Println("RPC Type is correct!")
			assertEqual(t, rt.Me.ID.Equals(&sentMessage.SenderId), true)
			fmt.Println("Sender ID is correct!")

			returnMessage := rpc.Message{rpc.TIME_OUT, *doNotCareID, *doNotCareID, []byte{byte(0)}}

			mbList := messageBufferList.GetInstance()
			buffer, _ := mbList.GetMessageBuffer(&sentMessage.RpcId)
			go buffer.AppendMessage(&returnMessage)

		},
	})

	net := testNetwork{}
	net.CheckList = cList
	kadem.SetNetworkHandler(&net)

	kadem.addContact(init("F000000000000000000000000000000000000000", "10.10.10.10:1000"))
	kadem.addContact(init("F000000000000000000000000000000000000001", "10.10.10.11:1000"))
	kadem.addContact(init("F000000000000000000000000000000000000002", "10.10.10.12:1000"))
	kadem.addContact(init("F000000000000000000000000000000000000003", "10.10.10.13:1000"))
	kadem.addContact(init("F000000000000000000000000000000000000004", "10.10.10.14:1000"))
	kadem.addContact(init("F000000000000000000000000000000000000005", "10.10.10.15:1000"))
	kadem.addContact(init("F000000000000000000000000000000000000006", "10.10.10.16:1000"))
	kadem.addContact(init("F000000000000000000000000000000000000007", "10.10.10.17:1000"))
	kadem.addContact(init("F000000000000000000000000000000000000008", "10.10.10.18:1000"))
	kadem.addContact(init("F000000000000000000000000000000000000009", "10.10.10.19:1000"))
	kadem.addContact(init("F00000000000000000000000000000000000000A", "10.10.10.20:1000"))
	kadem.addContact(init("F00000000000000000000000000000000000000B", "10.10.10.21:1000"))
	kadem.addContact(init("F00000000000000000000000000000000000000C", "10.10.10.22:1000"))
	kadem.addContact(init("F00000000000000000000000000000000000000D", "10.10.10.23:1000"))
	kadem.addContact(init("F00000000000000000000000000000000000000E", "10.10.10.24:1000"))
	kadem.addContact(init("F00000000000000000000000000000000000000F", "10.10.10.25:1000"))
	kadem.addContact(init("F000000000000000000000000000000000000010", "10.10.10.26:1000"))
	kadem.addContact(init("F000000000000000000000000000000000000012", "10.10.10.27:1000"))
	kadem.addContact(init("F000000000000000000000000000000000000013", "10.10.10.28:1000"))
	kadem.addContact(init("F000000000000000000000000000000000000014", "10.10.10.29:1000"))

	//should be discarded
	kadem.addContact(init("F000000000000000000000000000000000000015", "10.10.10.30:1000"))

	//should replace 10.10.10.11
	kadem.addContact(init("F000000000000000000000000000000000000016", "10.10.10.31:1000"))
	checkdata["F000000000000000000000000000000000000001"] = false
	checkdata["10.10.10.11:1000"] = false
	checkdata["F000000000000000000000000000000000000016"] = true
	checkdata["10.10.10.31:1000"] = true

	contacts := rt.FindClosestContacts(d7024e.NewKademliaID("F000000000000000000000000000000000000000"), 20)
	nrC := 20
	for _, c := range contacts {
		assertEqual(t, checkdata[strings.ToUpper(c.ID.String())], true)
		assertEqual(t, checkdata[c.Address], true)
		checkdata[strings.ToUpper(c.ID.String())] = false
		checkdata[c.Address] = false
		nrC--
	}
	fmt.Println("All correct contacts was found! ")
	assertEqual(t, nrC, 0)
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
