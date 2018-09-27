package kademlia

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"

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

	go checkData.CheckFunction(sentMessage, addr)

}

func (net *testNetwork) FetchFile(a string, b string) error {
	return nil
}

func TestFindNode(t *testing.T) {
	target := d7024e.NewContact(d7024e.NewKademliaID("FFFFFFFFF0000000000000000000000000000000"), "localhost:8000")
	kadem := GetInstance()

	startingContacts := []d7024e.Contact{
		d7024e.NewContact(d7024e.NewKademliaID("FFFFFFFFF0000000000000000000000000000001"), "localhost:8001"),
		d7024e.NewContact(d7024e.NewKademliaID("FFFFFFFFF0000000000000000000000000000002"), "localhost:8002"),
		d7024e.NewContact(d7024e.NewKademliaID("FFFFFFFFF0000000000000000000000000000003"), "localhost:8003"),
		d7024e.NewContact(d7024e.NewKademliaID("FFFFFFFFF0000000000000000000000000000004"), "localhost:8004"),
		d7024e.NewContact(d7024e.NewKademliaID("FFFFFFFFF0000000000000000000000000000005"), "localhost:8005"),
		d7024e.NewContact(d7024e.NewKademliaID("FFFFFFFFF0000000000000000000000000000006"), "localhost:8006"),
		d7024e.NewContact(d7024e.NewKademliaID("FFFFFFFFF0000000000000000000000000000015"), "localhost:8021"),
		d7024e.NewContact(d7024e.NewKademliaID("FFFFFFFFF0000000000000000000000000000016"), "localhost:8022"),
		d7024e.NewContact(d7024e.NewKademliaID("FFFFFFFFF0000000000000000000000000000017"), "localhost:8023"),
		d7024e.NewContact(d7024e.NewKademliaID("FFFFFFFFF0000000000000000000000000000018"), "localhost:8024"),
		d7024e.NewContact(d7024e.NewKademliaID("FFFFFFFFF0000000000000000000000000000019"), "localhost:8025"),
		d7024e.NewContact(d7024e.NewKademliaID("FFFFFFFFF000000000000000000000000000001A"), "localhost:8026"),
		d7024e.NewContact(d7024e.NewKademliaID("FFFFFFFFF000000000000000000000000000001B"), "localhost:8027"),
		d7024e.NewContact(d7024e.NewKademliaID("FFFFFFFFF000000000000000000000000000001C"), "localhost:8028"),
		d7024e.NewContact(d7024e.NewKademliaID("FFFFFFFFF000000000000000000000000000001D"), "localhost:8029"),
	}

	recipients := []d7024e.Contact{
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
		d7024e.NewContact(d7024e.NewKademliaID("FFFFFFFFF0000000000000000000000000000015"), "localhost:8021"),
		d7024e.NewContact(d7024e.NewKademliaID("FFFFFFFFF0000000000000000000000000000016"), "localhost:8022"),
		d7024e.NewContact(d7024e.NewKademliaID("FFFFFFFFF0000000000000000000000000000017"), "localhost:8023"),
	}


	expectedResult := []d7024e.Contact{
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
	
	recipientIndexLock := sync.Mutex
	recipientIndex := 0

	cList := []testNetworkControl{
		testNetworkControl{ //First reponse
			func(msg rpc.Message, addr string) {

				recipientIndexLock.Lock()
				expectedRecipitent := recipients[recipientIndex]
				expectedType := rpc.FIND_NODE
				fmt.Println("expected ", expectedRecipitent.Address, " got ", addr)
				assertEqual(t, expectedType, msg.RpcType)
				assertEqual(t, expectedRecipitent.Address, addr)

				nodesFound := []d7024e.Contact{
					d7024e.NewContact(d7024e.NewKademliaID("FFFFFFFFF0000000000000000000000000000007"), "localhost:8007"),
					d7024e.NewContact(d7024e.NewKademliaID("FFFFFFFFF0000000000000000000000000000008"), "localhost:8008"),
					d7024e.NewContact(d7024e.NewKademliaID("FFFFFFFFF0000000000000000000000000000009"), "localhost:8009"),
					d7024e.NewContact(d7024e.NewKademliaID("FFFFFFFFF000000000000000000000000000000A"), "localhost:8010"),
					d7024e.NewContact(d7024e.NewKademliaID("FFFFFFFFF000000000000000000000000000000B"), "localhost:8011"),
					d7024e.NewContact(d7024e.NewKademliaID("FFFFFFFFF000000000000000000000000000000C"), "localhost:8012"),
				}

				for _, c := range nodesFound {
					c.CalcDistance(target.ID)
				}

				nodesFoundM, _ := json.Marshal(nodesFound)
				firstResponse := rpc.Message{rpc.CLOSEST_NODES, msg.RpcId, *recipients[recipientIndex].ID, nodesFoundM}
				byteMsg, _ := json.Marshal(firstResponse)
				recipientIndex++
				recipientIndexLock.Unlock()
				kadem.HandleIncomingRPC(byteMsg, addr)

			},
		},

		testNetworkControl{ //Second reponse
			func(msg rpc.Message, addr string) {

				recipientIndexLock.Lock()
				expectedRecipitent := recipients[recipientIndex]
				expectedType := rpc.FIND_NODE
				fmt.Println("expected ", expectedRecipitent.Address, " got ", addr)
				assertEqual(t, expectedType, msg.RpcType)
				assertEqual(t, expectedRecipitent.Address, addr)

				nodesFound := []d7024e.Contact{
					d7024e.NewContact(d7024e.NewKademliaID("FFFFFFFFF000000000000000000000000000000D"), "localhost:8013"),
					d7024e.NewContact(d7024e.NewKademliaID("FFFFFFFFF000000000000000000000000000000E"), "localhost:8014"),
					d7024e.NewContact(d7024e.NewKademliaID("FFFFFFFFF000000000000000000000000000000F"), "localhost:8015"),
					d7024e.NewContact(d7024e.NewKademliaID("FFFFFFFFF0000000000000000000000000000010"), "localhost:8016"),
					d7024e.NewContact(d7024e.NewKademliaID("FFFFFFFFF0000000000000000000000000000011"), "localhost:8017"),
					d7024e.NewContact(d7024e.NewKademliaID("FFFFFFFFF0000000000000000000000000000012"), "localhost:8018"),
				}

				for _, c := range nodesFound {
					c.CalcDistance(target.ID)
				}

				nodesFoundM, _ := json.Marshal(nodesFound)
				secondResponse := rpc.Message{rpc.CLOSEST_NODES, msg.RpcId, *recipients[recipientIndex].ID, nodesFoundM}
				byteMsg2, _ := json.Marshal(secondResponse)
				recipientIndex++
				recipientIndexLock.Unlock()
				kadem.HandleIncomingRPC(byteMsg2, addr)
				fmt.Println("ejo2")
			},
		},

		testNetworkControl{ //Third reponse
			func(msg rpc.Message, addr string) {
				
				recipientIndexLock.Lock()
				expectedRecipitent := recipients[recipientIndex]
				expectedType := rpc.FIND_NODE
				fmt.Println("expected ", expectedRecipitent.Address, " got ", addr)
				assertEqual(t, expectedType, msg.RpcType)
				assertEqual(t, expectedRecipitent.Address, addr)

				nodesFound := []d7024e.Contact{
					d7024e.NewContact(d7024e.NewKademliaID("FFFFFFFFF0000000000000000000000000000013"), "localhost:8019"),
					d7024e.NewContact(d7024e.NewKademliaID("FFFFFFFFF0000000000000000000000000000014"), "localhost:8020"),
				}
				for _, c := range nodesFound {
					c.CalcDistance(target.ID)
				}

				nodesFoundM, _ := json.Marshal(nodesFound)
				thirdResponse := rpc.Message{rpc.CLOSEST_NODES, msg.RpcId, *recipients[recipientIndex].ID, nodesFoundM}
				byteMsg, _ := json.Marshal(thirdResponse)
				recipientIndex++
				recipientIndexLock.Unlock()

				kadem.HandleIncomingRPC(byteMsg, addr)
				fmt.Println("ejo3")
			},
		},
	}

	//Rest of the responses
	for i := 0; i < 20; i++ {
		cList = append(cList, testNetworkControl{
			func(msg rpc.Message, addr string) {
				recipientIndexLock.Lock()
				expectedRecipitent := recipients[recipientIndex]
				expectedType := rpc.FIND_NODE
				fmt.Println("expected ", expectedRecipitent.Address, " got ", addr)
				assertEqual(t, expectedType, msg.RpcType)
				assertEqual(t, expectedRecipitent.Address, addr)

				nodesFound := []d7024e.Contact{
					d7024e.NewContact(d7024e.NewKademliaID("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF"), "localhost:8021"),
				}
				nodesFound[0].CalcDistance(target.ID)

				nodesFoundM, _ := json.Marshal(nodesFound)
				response := rpc.Message{rpc.CLOSEST_NODES, msg.RpcId, *recipients[recipientIndex].ID, nodesFoundM}
				byteMsg, _ := json.Marshal(response)
				recipientIndex++
				recipientIndexLock.Unlock()
				
				kadem.HandleIncomingRPC(byteMsg, addr)
			},
		})
	}

	//Routing table setup
	rt := routingTable.GetInstance()
	for _, c := range startingContacts {
		rt.AddContact(c)
	}

	net := testNetwork{}
	net.CheckList = cList
	kadem.SetNetworkHandler(&net)

	result, _, _ := kadem.lookupProcedure(procedureContacts, target.ID)

	assertEqual(t, expectedResult, result)
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

				GetInstance().HandleIncomingRPC(d, "10.10.10.10:1000")

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
	var wg1, wg2 sync.WaitGroup
	wg1.Add(1)
	wg2.Add(1)

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

			GetInstance().HandleIncomingRPC(d, addr)
			wg1.Done()
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
			buffer.AppendMessage(&returnMessage)
			time.Sleep(2 * time.Second)
			wg2.Done()
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
	wg1.Wait()

	//should replace 10.10.10.11
	kadem.addContact(init("F000000000000000000000000000000000000016", "10.10.10.31:1000"))
	wg2.Wait()

	checkdata["F000000000000000000000000000000000000001"] = false
	checkdata["10.10.10.11:1000"] = false
	checkdata["F000000000000000000000000000000000000016"] = true
	checkdata["10.10.10.31:1000"] = true

	contacts := rt.FindClosestContacts(d7024e.NewKademliaID("F000000000000000000000000000000000000000"), 20)
	nrC := 20
	for _, c := range contacts {
		fmt.Println("testing ", c.Address)
		assertEqual(t, checkdata[strings.ToUpper(c.ID.String())], true)
		assertEqual(t, checkdata[c.Address], true)
		checkdata[strings.ToUpper(c.ID.String())] = false
		checkdata[c.Address] = false
		nrC--
	}
	fmt.Println("All correct contacts was found! ")
	assertEqual(t, nrC, 0)
}

func TestPing(t *testing.T) {
	rt := routingTable.GetInstance()
	kadem := GetInstance()
	var wg1 sync.WaitGroup
	wg1.Add(1)

	RpcId := d7024e.NewRandomKademliaID()
	SendId := d7024e.NewRandomKademliaID()

	var cList []testNetworkControl
	cList = append(cList, testNetworkControl{
		func(sentMessage rpc.Message, addr string) {
			assertEqual(t, addr, "10.10.10.10:1000")
			fmt.Println("correct address is responded to!")

			assertEqual(t, sentMessage.RpcType, rpc.PONG)
			fmt.Println("RPC Type is correct!")
			assertEqual(t, rt.Me.ID.Equals(&sentMessage.SenderId), true)
			fmt.Println("Sender ID is correct!")
			assertEqual(t, RpcId.Equals(&sentMessage.RpcId), true)
			fmt.Println("Correct RPCID is used!")

			wg1.Done()
		},
	})

	msg := rpc.Message{rpc.PING, *RpcId, *SendId, []byte{byte(0)}}
	d, _ := json.Marshal(msg)

	net := testNetwork{}
	net.CheckList = cList
	kadem.SetNetworkHandler(&net)

	kadem.HandleIncomingRPC(d, "10.10.10.10:1000")
	wg1.Wait()
}

func helperReturnMarshal(data interface{}) []byte {
	da, _ := json.Marshal(data)
	return da
}

func assertEqual(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		panic(fmt.Sprintf("%s != %s", a, b))
		//t.Fatalf("%s != %s", a, b)
	}
}
