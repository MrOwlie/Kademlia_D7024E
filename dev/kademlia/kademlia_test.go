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
	"../metadata"
	"../routingTable"
	"../rpc"
)

var doNotCareID *d7024e.KademliaID = d7024e.NewKademliaID("F000000000000000000000000000000000000000")

type testNetworkControl struct {
	CheckFunction func(rpc.Message, string)
	//ReturnMsg     rpc.Message
	//FromAddress   string
}

type testNetworkFetchControl struct {
	CheckFunction func(string, string)
	//ReturnMsg     rpc.Message
	//FromAddress   string
}

type testNetwork struct {
	CheckList      []testNetworkControl
	CheckListFetch []testNetworkFetchControl
	sendMutex      sync.Mutex
	fetchMutex     sync.Mutex
}

func (net *testNetwork) SendMessage(addr string, data *[]byte) {
	net.sendMutex.Lock()
	defer net.sendMutex.Unlock()
	fmt.Printf("Sending message to :%v\n", addr)
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

func (net *testNetwork) FetchFile(url string, filePath string) error {

	net.fetchMutex.Lock()
	defer net.fetchMutex.Unlock()
	fmt.Printf("Fetching file from :%v\n", url)
	checkData := net.CheckListFetch[0]
	if len(net.CheckListFetch) > 1 {
		net.CheckListFetch = net.CheckListFetch[1:]
	} else {
		net.CheckListFetch = nil
	}

	go checkData.CheckFunction(url, filePath)

	return nil
}

type expectedRecipient struct {
	valid       bool
	recipientId *d7024e.KademliaID
}

var firstAlphaRecipients map[string]*expectedRecipient = make(map[string]*expectedRecipient)
var finalKRecipients map[string]*expectedRecipient = make(map[string]*expectedRecipient)

//Base test for LookupContact
func lookUpTestInitalSetup(target d7024e.Contact, t *testing.T, kadem *Kademlia, rt *routingTable.RoutingTable) *[]testNetworkControl {

	startingContacts := []d7024e.Contact{
		d7024e.NewContact(d7024e.NewKademliaID("FFFFFFFFF0000000000000000000000000000001"), "localhost:8001"),
		d7024e.NewContact(d7024e.NewKademliaID("FFFFFFFFF0000000000000000000000000000002"), "localhost:8002"),
		d7024e.NewContact(d7024e.NewKademliaID("FFFFFFFFF0000000000000000000000000000003"), "localhost:8003"),
		d7024e.NewContact(d7024e.NewKademliaID("FFFFFFFFF0000000000000000000000000000004"), "localhost:8004"),
		d7024e.NewContact(d7024e.NewKademliaID("FFFFFFFFF0000000000000000000000000000005"), "localhost:8005"),
		d7024e.NewContact(d7024e.NewKademliaID("FFFFFFFFF0000000000000000000000000000006"), "localhost:8006"),
		d7024e.NewContact(d7024e.NewKademliaID("0FFFFFFFF0000000000000000000000000000015"), "localhost:8021"),
		d7024e.NewContact(d7024e.NewKademliaID("0FFFFFFFF0000000000000000000000000000016"), "localhost:8022"),
		d7024e.NewContact(d7024e.NewKademliaID("0FFFFFFFF0000000000000000000000000000017"), "localhost:8023"),
		d7024e.NewContact(d7024e.NewKademliaID("0FFFFFFFF0000000000000000000000000000018"), "localhost:8024"),
		d7024e.NewContact(d7024e.NewKademliaID("0FFFFFFFF0000000000000000000000000000019"), "localhost:8025"),
		d7024e.NewContact(d7024e.NewKademliaID("0FFFFFFFF000000000000000000000000000001A"), "localhost:8026"),
		d7024e.NewContact(d7024e.NewKademliaID("0FFFFFFFF000000000000000000000000000001B"), "localhost:8027"),
		d7024e.NewContact(d7024e.NewKademliaID("0FFFFFFFF000000000000000000000000000001C"), "localhost:8028"),
		d7024e.NewContact(d7024e.NewKademliaID("0FFFFFFFF000000000000000000000000000001D"), "localhost:8029"),
	}

	firstAlphaRecipients["localhost:8001"] = &expectedRecipient{true, d7024e.NewKademliaID("FFFFFFFFF0000000000000000000000000000001")}
	firstAlphaRecipients["localhost:8002"] = &expectedRecipient{true, d7024e.NewKademliaID("FFFFFFFFF0000000000000000000000000000002")}
	firstAlphaRecipients["localhost:8003"] = &expectedRecipient{true, d7024e.NewKademliaID("FFFFFFFFF0000000000000000000000000000003")}

	finalKRecipients["localhost:8004"] = &expectedRecipient{true, d7024e.NewKademliaID("FFFFFFFFF0000000000000000000000000000004")}
	finalKRecipients["localhost:8005"] = &expectedRecipient{true, d7024e.NewKademliaID("FFFFFFFFF0000000000000000000000000000005")}
	finalKRecipients["localhost:8006"] = &expectedRecipient{true, d7024e.NewKademliaID("FFFFFFFFF0000000000000000000000000000006")}
	finalKRecipients["localhost:8007"] = &expectedRecipient{true, d7024e.NewKademliaID("FFFFFFFFF0000000000000000000000000000007")}
	finalKRecipients["localhost:8008"] = &expectedRecipient{true, d7024e.NewKademliaID("FFFFFFFFF0000000000000000000000000000008")}
	finalKRecipients["localhost:8009"] = &expectedRecipient{true, d7024e.NewKademliaID("FFFFFFFFF0000000000000000000000000000009")}
	finalKRecipients["localhost:8010"] = &expectedRecipient{true, d7024e.NewKademliaID("FFFFFFFFF000000000000000000000000000000A")}
	finalKRecipients["localhost:8011"] = &expectedRecipient{true, d7024e.NewKademliaID("FFFFFFFFF000000000000000000000000000000B")}
	finalKRecipients["localhost:8012"] = &expectedRecipient{true, d7024e.NewKademliaID("FFFFFFFFF000000000000000000000000000000C")}
	finalKRecipients["localhost:8013"] = &expectedRecipient{true, d7024e.NewKademliaID("FFFFFFFFF000000000000000000000000000000D")}
	finalKRecipients["localhost:8014"] = &expectedRecipient{true, d7024e.NewKademliaID("FFFFFFFFF000000000000000000000000000000E")}
	finalKRecipients["localhost:8015"] = &expectedRecipient{true, d7024e.NewKademliaID("FFFFFFFFF000000000000000000000000000000F")}
	finalKRecipients["localhost:8016"] = &expectedRecipient{true, d7024e.NewKademliaID("FFFFFFFFF0000000000000000000000000000010")}
	finalKRecipients["localhost:8017"] = &expectedRecipient{true, d7024e.NewKademliaID("FFFFFFFFF0000000000000000000000000000011")}
	finalKRecipients["localhost:8018"] = &expectedRecipient{true, d7024e.NewKademliaID("FFFFFFFFF0000000000000000000000000000012")}
	finalKRecipients["localhost:8019"] = &expectedRecipient{true, d7024e.NewKademliaID("FFFFFFFFF0000000000000000000000000000013")}
	finalKRecipients["localhost:8020"] = &expectedRecipient{true, d7024e.NewKademliaID("FFFFFFFFF0000000000000000000000000000014")}
	finalKRecipients["localhost:8021"] = &expectedRecipient{true, d7024e.NewKademliaID("0FFFFFFFF0000000000000000000000000000015")}
	finalKRecipients["localhost:8022"] = &expectedRecipient{true, d7024e.NewKademliaID("0FFFFFFFF0000000000000000000000000000016")}
	finalKRecipients["localhost:8023"] = &expectedRecipient{true, d7024e.NewKademliaID("0FFFFFFFF0000000000000000000000000000017")}

	cList := []testNetworkControl{
		testNetworkControl{ //First reponse
			func(msg rpc.Message, addr string) {

				if recipient, ok := firstAlphaRecipients[addr]; ok && recipient.valid {
					recipient.valid = false
				} else {
					t.Fail()
				}

				expectedType := rpc.FIND_NODE
				assertEqual(t, expectedType, msg.RpcType)

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

				nodesFoundM, _ := json.Marshal(rpc.ClosestNodes{nodesFound})
				firstResponse := rpc.Message{rpc.CLOSEST_NODES, msg.RpcId, *firstAlphaRecipients[addr].recipientId, nodesFoundM}
				byteMsg, _ := json.Marshal(firstResponse)

				kadem.HandleIncomingRPC(byteMsg, addr)

			},
		},

		testNetworkControl{ //Second reponse
			func(msg rpc.Message, addr string) {

				if recipient, ok := firstAlphaRecipients[addr]; ok && recipient.valid {
					recipient.valid = false
				} else {
					t.Fail()
				}

				expectedType := rpc.FIND_NODE
				assertEqual(t, expectedType, msg.RpcType)

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

				nodesFoundM, _ := json.Marshal(rpc.ClosestNodes{nodesFound})
				secondResponse := rpc.Message{rpc.CLOSEST_NODES, msg.RpcId, *firstAlphaRecipients[addr].recipientId, nodesFoundM}
				byteMsg2, _ := json.Marshal(secondResponse)

				kadem.HandleIncomingRPC(byteMsg2, addr)
				fmt.Println("ejo2")
			},
		},

		testNetworkControl{ //Third reponse
			func(msg rpc.Message, addr string) {

				if recipient, ok := firstAlphaRecipients[addr]; ok && recipient.valid {
					recipient.valid = false
				} else {
					t.Fail()
				}

				expectedType := rpc.FIND_NODE
				assertEqual(t, expectedType, msg.RpcType)

				nodesFound := []d7024e.Contact{
					d7024e.NewContact(d7024e.NewKademliaID("FFFFFFFFF0000000000000000000000000000013"), "localhost:8019"),
					d7024e.NewContact(d7024e.NewKademliaID("FFFFFFFFF0000000000000000000000000000014"), "localhost:8020"),
				}
				for _, c := range nodesFound {
					c.CalcDistance(target.ID)
				}

				nodesFoundM, _ := json.Marshal(rpc.ClosestNodes{nodesFound})
				thirdResponse := rpc.Message{rpc.CLOSEST_NODES, msg.RpcId, *firstAlphaRecipients[addr].recipientId, nodesFoundM}
				byteMsg, _ := json.Marshal(thirdResponse)

				kadem.HandleIncomingRPC(byteMsg, addr)
				fmt.Println("ejo3")
			},
		},
	}

	//Rest of the responses
	for i := 0; i < 20; i++ {
		cList = append(cList, testNetworkControl{
			func(msg rpc.Message, addr string) {

				if recipient, ok := finalKRecipients[addr]; ok && recipient.valid {
					recipient.valid = false
				} else {
					fmt.Printf("No message to %v expected", addr)
					//t.FailNow()
				}

				expectedType := rpc.FIND_NODE
				assertEqual(t, expectedType, msg.RpcType)

				nodesFound := []d7024e.Contact{
					d7024e.NewContact(d7024e.NewKademliaID("0FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF"), "localhost:8090"),
				}
				nodesFound[0].CalcDistance(target.ID)

				nodesFoundM, _ := json.Marshal(rpc.ClosestNodes{nodesFound})
				response := rpc.Message{rpc.CLOSEST_NODES, msg.RpcId, *finalKRecipients[addr].recipientId, nodesFoundM}
				byteMsg, _ := json.Marshal(response)
				kadem.HandleIncomingRPC(byteMsg, addr)
			},
		})
	}

	//Routing table setup
	for _, c := range startingContacts {
		rt.AddContact(c)
	}

	return &cList
}

func TestFindNode(t *testing.T) {
	target := d7024e.NewContact(d7024e.NewKademliaID("FFFFFFFFF0000000000000000000000000000000"), "localhost:8000")
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

	RTable := routingTable.NewRoutingTable()
	MBList := &messageBufferList.MessageBufferList{}
	MData := metadata.NewFileMetaData()
	var kadem = NewKademliaObject(RTable, MBList, MData)
	net := testNetwork{}
	net.CheckList = *lookUpTestInitalSetup(target, t, kadem, RTable)
	kadem.SetNetworkHandler(&net)

	result, _, _ := kadem.lookupProcedure(procedureContacts, target.ID)
	if len(result) == len(expectedResult) {
		for i, c := range result {
			assertEqual(t, c.ID.String(), expectedResult[i].ID.String())
		}
	} else {
		t.Fail()
	}

}

func TestFindNodeTimeOut(t *testing.T) {
	target := d7024e.NewContact(d7024e.NewKademliaID("FFFFFFFFF0000000000000000000000000000000"), "localhost:8000")
	expectedResult := []d7024e.Contact{
		d7024e.NewContact(d7024e.NewKademliaID("FFFFFFFFF0000000000000000000000000000001"), "localhost:8001"),
		d7024e.NewContact(d7024e.NewKademliaID("FFFFFFFFF0000000000000000000000000000002"), "localhost:8002"),
		d7024e.NewContact(d7024e.NewKademliaID("FFFFFFFFF0000000000000000000000000000003"), "localhost:8003"),
		d7024e.NewContact(d7024e.NewKademliaID("FFFFFFFFF0000000000000000000000000000004"), "localhost:8004"),
		d7024e.NewContact(d7024e.NewKademliaID("FFFFFFFFF0000000000000000000000000000005"), "localhost:8005"),
		d7024e.NewContact(d7024e.NewKademliaID("FFFFFFFFF0000000000000000000000000000006"), "localhost:8006"),
		d7024e.NewContact(d7024e.NewKademliaID("FFFFFFFFF000000000000000000000000000000D"), "localhost:8013"),
		d7024e.NewContact(d7024e.NewKademliaID("FFFFFFFFF000000000000000000000000000000E"), "localhost:8014"),
		d7024e.NewContact(d7024e.NewKademliaID("FFFFFFFFF000000000000000000000000000000F"), "localhost:8015"),
		d7024e.NewContact(d7024e.NewKademliaID("FFFFFFFFF0000000000000000000000000000010"), "localhost:8016"),
		d7024e.NewContact(d7024e.NewKademliaID("FFFFFFFFF0000000000000000000000000000011"), "localhost:8017"),
		d7024e.NewContact(d7024e.NewKademliaID("FFFFFFFFF0000000000000000000000000000012"), "localhost:8018"),
		d7024e.NewContact(d7024e.NewKademliaID("FFFFFFFFF0000000000000000000000000000013"), "localhost:8019"),
		d7024e.NewContact(d7024e.NewKademliaID("FFFFFFFFF0000000000000000000000000000014"), "localhost:8020"),
		d7024e.NewContact(d7024e.NewKademliaID("0FFFFFFFF0000000000000000000000000000015"), "localhost:8021"),
		d7024e.NewContact(d7024e.NewKademliaID("0FFFFFFFF0000000000000000000000000000016"), "localhost:8022"),
		d7024e.NewContact(d7024e.NewKademliaID("0FFFFFFFF0000000000000000000000000000017"), "localhost:8023"),
		d7024e.NewContact(d7024e.NewKademliaID("0FFFFFFFF0000000000000000000000000000018"), "localhost:8024"),
		d7024e.NewContact(d7024e.NewKademliaID("0FFFFFFFF0000000000000000000000000000019"), "localhost:8025"),
		d7024e.NewContact(d7024e.NewKademliaID("0FFFFFFFF000000000000000000000000000001A"), "localhost:8026"),
		d7024e.NewContact(d7024e.NewKademliaID("0FFFFFFFF000000000000000000000000000001B"), "localhost:8027"),
	}

	RTable := routingTable.NewRoutingTable()
	MBList := &messageBufferList.MessageBufferList{}
	MData := metadata.NewFileMetaData()
	var kadem = NewKademliaObject(RTable, MBList, MData)

	cList := *lookUpTestInitalSetup(target, t, kadem, RTable)
	cList[0] = testNetworkControl{func(msg rpc.Message, addr string) {
		if recipient, ok := firstAlphaRecipients[addr]; ok && recipient.valid {
			recipient.valid = false
		} else {
			fmt.Printf("No message to %v expected.", addr)
			t.Fail()
		}

		if expectedResult[0].Address == addr{
			expectedResult = expectedResult[1:]
		} else if expectedResult[1].Address == addr {
			expectedResult = append(expectedResult[:1], expectedResult[2:]...)
		} else if expectedResult[2].Address == addr {
			expectedResult = append(expectedResult[:2], expectedResult[3:]...)
		} else {
			fmt.Printf("No message to %v expected", addr)
			t.Fail()
		}
	}}

	finalKRecipients["localhost:8007"].valid = false
	finalKRecipients["localhost:8008"].valid = false
	finalKRecipients["localhost:8009"].valid = false
	finalKRecipients["localhost:8010"].valid = false
	finalKRecipients["localhost:8011"].valid = false
	finalKRecipients["localhost:8012"].valid = false

	finalKRecipients["localhost:8024"] = &expectedRecipient{true, d7024e.NewKademliaID("0FFFFFFFF0000000000000000000000000000018")}
	finalKRecipients["localhost:8025"] = &expectedRecipient{true, d7024e.NewKademliaID("0FFFFFFFF0000000000000000000000000000019")}
	finalKRecipients["localhost:8026"] = &expectedRecipient{true, d7024e.NewKademliaID("0FFFFFFFF000000000000000000000000000001A")}
	finalKRecipients["localhost:8027"] = &expectedRecipient{true, d7024e.NewKademliaID("0FFFFFFFF000000000000000000000000000001B")}
	finalKRecipients["localhost:8028"] = &expectedRecipient{true, d7024e.NewKademliaID("0FFFFFFFF000000000000000000000000000001C")}
	finalKRecipients["localhost:8029"] = &expectedRecipient{true, d7024e.NewKademliaID("0FFFFFFFF000000000000000000000000000001D")}

	net := testNetwork{}
	net.CheckList = cList
	kadem.SetNetworkHandler(&net)

	result, _, _ := kadem.lookupProcedure(procedureContacts, target.ID)
	if len(result) == len(expectedResult) {
		for i, c := range result {
			assertEqual(t, c.ID.String(), expectedResult[i].ID.String())
		}
	} else {
		panic(fmt.Sprintf("Unexpected amount of contacts: %v\n", len(result)))
	}
}

func TestFindNodeTimeOutRecovery(t *testing.T) {
	target := d7024e.NewContact(d7024e.NewKademliaID("FFFFFFFFF0000000000000000000000000000000"), "localhost:8000")
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

	RTable := routingTable.NewRoutingTable()
	MBList := &messageBufferList.MessageBufferList{}
	MData := metadata.NewFileMetaData()
	var kadem = NewKademliaObject(RTable, MBList, MData)

	cList := *lookUpTestInitalSetup(target, t, kadem, RTable)

	cList[0] = testNetworkControl{ //Third reponse
		func(msg rpc.Message, addr string) {

			if recipient, ok := firstAlphaRecipients[addr]; ok && recipient.valid {
				recipient.valid = false
			} else {
				t.Fail()
			}

			expectedType := rpc.FIND_NODE
			assertEqual(t, expectedType, msg.RpcType)

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

			nodesFoundM, _ := json.Marshal(rpc.ClosestNodes{nodesFound})
			firstResponse := rpc.Message{rpc.CLOSEST_NODES, msg.RpcId, *firstAlphaRecipients[addr].recipientId, nodesFoundM}
			byteMsg, _ := json.Marshal(firstResponse)
			time.Sleep(1 * time.Second)
			kadem.HandleIncomingRPC(byteMsg, addr)
			fmt.Println("ejo3")
		},
	}

	cList[22] = testNetworkControl{
		func(msg rpc.Message, addr string) {

			if recipient, ok := finalKRecipients[addr]; ok && recipient.valid {
				recipient.valid = false
			} else {
				fmt.Printf("No message to %v expected", addr)
				//t.FailNow()
			}

			expectedType := rpc.FIND_NODE
			assertEqual(t, expectedType, msg.RpcType)

			nodesFound := []d7024e.Contact{
				d7024e.NewContact(d7024e.NewKademliaID("0FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF"), "localhost:8090"),
			}
			nodesFound[0].CalcDistance(target.ID)

			nodesFoundM, _ := json.Marshal(rpc.ClosestNodes{nodesFound})
			response := rpc.Message{rpc.CLOSEST_NODES, msg.RpcId, *finalKRecipients[addr].recipientId, nodesFoundM}
			byteMsg, _ := json.Marshal(response)
			time.Sleep(2 * time.Second)
			kadem.HandleIncomingRPC(byteMsg, addr)
		},
	}

	finalKRecipients["localhost:8007"].valid = false
	finalKRecipients["localhost:8008"].valid = false
	finalKRecipients["localhost:8009"].valid = false
	finalKRecipients["localhost:8010"].valid = false
	finalKRecipients["localhost:8011"].valid = false
	finalKRecipients["localhost:8012"].valid = false

	finalKRecipients["localhost:8024"] = &expectedRecipient{true, d7024e.NewKademliaID("0FFFFFFFF0000000000000000000000000000018")}
	finalKRecipients["localhost:8025"] = &expectedRecipient{true, d7024e.NewKademliaID("0FFFFFFFF0000000000000000000000000000019")}
	finalKRecipients["localhost:8026"] = &expectedRecipient{true, d7024e.NewKademliaID("0FFFFFFFF000000000000000000000000000001A")}
	finalKRecipients["localhost:8027"] = &expectedRecipient{true, d7024e.NewKademliaID("0FFFFFFFF000000000000000000000000000001B")}
	finalKRecipients["localhost:8028"] = &expectedRecipient{true, d7024e.NewKademliaID("0FFFFFFFF000000000000000000000000000001C")}
	finalKRecipients["localhost:8029"] = &expectedRecipient{true, d7024e.NewKademliaID("0FFFFFFFF000000000000000000000000000001D")}

	net := testNetwork{}
	net.CheckList = cList
	kadem.SetNetworkHandler(&net)

	result, _, _ := kadem.lookupProcedure(procedureContacts, target.ID)
	if len(result) == len(expectedResult) {
		for i, c := range result {
			assertEqual(t, c.ID.String(), expectedResult[i].ID.String())
		}
	} else {
		t.Fail()
	}
}

func TestFindValue(t *testing.T) {

	RTable := routingTable.NewRoutingTable()
	MBList := &messageBufferList.MessageBufferList{}
	MData := metadata.NewFileMetaData()
	var kadem = NewKademliaObject(RTable, MBList, MData)

	target := d7024e.NewContact(d7024e.NewKademliaID("FFFFFFFFF0000000000000000000000000000000"), "localhost:8030")

	startingContacts := []d7024e.Contact{
		d7024e.NewContact(d7024e.NewKademliaID("FFFFFFFFF0000000000000000000000000000001"), "localhost:8001"),
		d7024e.NewContact(d7024e.NewKademliaID("FFFFFFFFF0000000000000000000000000000002"), "localhost:8002"),
		d7024e.NewContact(d7024e.NewKademliaID("FFFFFFFFF0000000000000000000000000000003"), "localhost:8003"),
	}

	firstAlphaRecipients["localhost:8001"] = &expectedRecipient{true, d7024e.NewKademliaID("FFFFFFFFF0000000000000000000000000000001")}
	firstAlphaRecipients["localhost:8002"] = &expectedRecipient{true, d7024e.NewKademliaID("FFFFFFFFF0000000000000000000000000000002")}
	firstAlphaRecipients["localhost:8003"] = &expectedRecipient{true, d7024e.NewKademliaID("FFFFFFFFF0000000000000000000000000000003")}

	var storeRecipient string

	urlTail := "/storage/FFFFFFFFF0000000000000000000000000000000"
	var fetchURL1 string
	var fetchURL2 string

	cList := []testNetworkControl{
		testNetworkControl{ //First reponse
			func(msg rpc.Message, addr string) {

				if recipient, ok := firstAlphaRecipients[addr]; ok && recipient.valid {
					recipient.valid = false
				} else {
					t.Fail()
				}

				expectedType := rpc.FIND_VALUE
				assertEqual(t, expectedType, msg.RpcType)

				firstResponse := rpc.Message{rpc.HAS_VALUE, msg.RpcId, *firstAlphaRecipients[addr].recipientId, []byte{byte(0)}}
				byteMsg, _ := json.Marshal(firstResponse)

				fetchURL1 = addr+urlTail

				kadem.HandleIncomingRPC(byteMsg, addr)
			},
		},

		testNetworkControl{ //Second reponse
			func(msg rpc.Message, addr string) {

				if recipient, ok := firstAlphaRecipients[addr]; ok && recipient.valid {
					recipient.valid = false
				} else {
					t.Fail()
				}

				expectedType := rpc.FIND_VALUE
				assertEqual(t, expectedType, msg.RpcType)

				firstResponse := rpc.Message{rpc.HAS_VALUE, msg.RpcId, *firstAlphaRecipients[addr].recipientId, []byte{byte(0)}}
				byteMsg, _ := json.Marshal(firstResponse)

				fetchURL2 = addr+urlTail

				kadem.HandleIncomingRPC(byteMsg, addr)
			},
		},

		testNetworkControl{ //Third reponse
			func(msg rpc.Message, addr string) {

				if recipient, ok := firstAlphaRecipients[addr]; ok && recipient.valid {
					recipient.valid = false
				} else {
					t.Fail()
				}

				storeRecipient = addr
				expectedType := rpc.FIND_VALUE
				assertEqual(t, expectedType, msg.RpcType)

				nodesFound := []d7024e.Contact{
					d7024e.NewContact(d7024e.NewKademliaID("FFFFFFFFF0000000000000000000000000000013"), "localhost:8019"),
					d7024e.NewContact(d7024e.NewKademliaID("FFFFFFFFF0000000000000000000000000000014"), "localhost:8020"),
				}
				for _, c := range nodesFound {
					c.CalcDistance(target.ID)
				}

				nodesFoundM, _ := json.Marshal(rpc.ClosestNodes{nodesFound})
				thirdResponse := rpc.Message{rpc.CLOSEST_NODES, msg.RpcId, *firstAlphaRecipients[addr].recipientId, nodesFoundM}
				byteMsg, _ := json.Marshal(thirdResponse)

				kadem.HandleIncomingRPC(byteMsg, addr)
				fmt.Println("ejo3")
			},
		},
		testNetworkControl{ //Cache in node
			func(msg rpc.Message, addr string) {

				if storeRecipient != addr {
					fmt.Println(addr)
					t.Fail()
				}

				expectedType := rpc.STORE
				assertEqual(t, expectedType, msg.RpcType)
			},
		},
	}

	fetchList := []testNetworkFetchControl{testNetworkFetchControl{func(url string, path string) {
		if url != fetchURL1 && url != fetchURL2 {
			t.Fail()
		}
	},
	},
	}

	rt := RTable
	for i := 0; i < d7024e.IDLength*8; i++ {
		rt.Buckets[i] = d7024e.NewBucket()
	}
	for _, c := range startingContacts {
		rt.AddContact(c)
	}

	net := testNetwork{}
	net.CheckList = cList
	net.CheckListFetch = fetchList

	kadem.SetNetworkHandler(&net)

	_, _, found := kadem.LookupData("FFFFFFFFF0000000000000000000000000000000")

	if !found {
		t.Fail()
	}
}

func TestJoin(t *testing.T) {
	RTable := routingTable.NewRoutingTable()
	MBList := &messageBufferList.MessageBufferList{}
	MData := metadata.NewFileMetaData()
	var kadem = NewKademliaObject(RTable, MBList, MData)

	rt := RTable
	for i := 0; i < d7024e.IDLength*8; i++ {
		rt.Buckets[i] = d7024e.NewBucket()
	}
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

				kadem.HandleIncomingRPC(d, "10.10.10.10:1000")

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
/*
func TestBucketReExploration(t *testing.T) {
	fmt.Println("STARTED")
	RTable := routingTable.NewRoutingTable()
	MBList := &messageBufferList.MessageBufferList{}
	MData := metadata.NewFileMetaData()
	var kadem = NewKademliaObject(RTable, MBList, MData)

	rt := RTable
	for i := 0; i < d7024e.IDLength*8; i++ {
		rt.Buckets[i] = d7024e.NewBucket()
	}
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
	fmt.Println("DONE")
}
*/
// func TestStoreFile(t *testing.T) {
// 	kadem := GetInstance()
// 	rt := routingTable.GetInstance()
// 	rt.Me.ID = d7024e.NewKademliaID("0000000000000000000000000000000000000000")
// 	rt.AddContact(d7024e.NewContact(d7024e.NewKademliaID("0000000000000000000000000000000000000001"), "localhost:8001"))
// 	//rt.AddContact(d7024e.NewContact(d7024e.NewKademliaID("0000000000000000000000000000000000000002"), "localhost:8002"))
// 	//rt.AddContact(d7024e.NewContact(d7024e.NewKademliaID("0000000000000000000000000000000000000003"), "localhost:8003"))
// 	//rt.AddContact(d7024e.NewContact(d7024e.NewKademliaID("0000000000000000000000000000000000000004"), "localhost:8004"))
// 	//rt.AddContact(d7024e.NewContact(d7024e.NewKademliaID("0000000000000000000000000000000000000005"), "localhost:8005"))
// 	//rt.AddContact(d7024e.NewContact(d7024e.NewKademliaID("0000000000000000000000000000000000000006"), "localhost:8006"))
// 	//rt.AddContact(d7024e.NewContact(d7024e.NewKademliaID("0000000000000000000000000000000000000007"), "localhost:8007"))
// 	//rt.AddContact(d7024e.NewContact(d7024e.NewKademliaID("0000000000000000000000000000000000000008"), "localhost:8008"))
// 	//rt.AddContact(d7024e.NewContact(d7024e.NewKademliaID("0000000000000000000000000000000000000009"), "localhost:8009"))

// 	downloadsPath, err := filepath.Abs("../downloads")
// 	storagePath, err := filepath.Abs("../storage")

// 	testContent := []byte("test")

// 	file, err := os.Create(downloadsPath + filepath.Clean("/test"))
// 	if err != nil {
// 		fmt.Println("TestStoreFile failed at creating test file in downloads.")
// 		fmt.Println(err)
// 	}

// 	file.Write(testContent)

// 	h := sha1.New()

// 	_, err = io.Copy(h, file)
// 	if err != nil {
// 		fmt.Println("TestStoreFile failed at copying test file in downloads.")
// 		fmt.Println(err)
// 	}
// 	file.Seek(0, 0)
// 	hash := h.Sum(nil)

// 	var wg sync.WaitGroup
// 	wg.Add(1)

// 	var cList []testNetworkControl
// 	cList = append(cList, testNetworkControl{
// 		func(sentMessage rpc.Message, addr string) {
// 			if sentMessage.RpcType != rpc.FIND_NODE {
// 				fmt.Println("TYPE: " + sentMessage.RpcType)
// 				panic("TestStoreFile failed because sent RPC type is not STORE.")
// 			} else {
// 				returnMessage := rpc.Message{rpc.CLOSEST_NODES, sentMessage.RpcId, *d7024e.NewKademliaID("0000000000000000000000000000000000000000"), helperReturnMarshal(rpc.ClosestNodes{
// 					[]d7024e.Contact{
// 						d7024e.NewContact(d7024e.NewKademliaID("0000000000000000000000000000000000000001"), "localhost:8001"),
// 						//d7024e.NewContact(d7024e.NewKademliaID("0000000000000000000000000000000000000002"), "localhost:8002"),
// 						//d7024e.NewContact(d7024e.NewKademliaID("0000000000000000000000000000000000000003"), "localhost:8003"),
// 						//d7024e.NewContact(d7024e.NewKademliaID("0000000000000000000000000000000000000004"), "localhost:8004"),
// 						//d7024e.NewContact(d7024e.NewKademliaID("0000000000000000000000000000000000000005"), "localhost:8005"),
// 						//d7024e.NewContact(d7024e.NewKademliaID("0000000000000000000000000000000000000006"), "localhost:8006"),
// 						//d7024e.NewContact(d7024e.NewKademliaID("0000000000000000000000000000000000000007"), "localhost:8007"),
// 						//d7024e.NewContact(d7024e.NewKademliaID("0000000000000000000000000000000000000008"), "localhost:8008"),
// 					},
// 				}),
// 				}
// 				d, _ := json.Marshal(returnMessage)

// 				GetInstance().HandleIncomingRPC(d, addr)
// 				wg.Done()
// 			}

// 		},
// 	})

// 	net := testNetwork{}
// 	net.CheckList = cList
// 	kadem.SetNetworkHandler(&net)
// 	kadem.StoreFile(downloadsPath + filepath.Clean("/test"))
// 	wg.Wait()

// 	_, errs := ioutil.ReadFile(storagePath + filepath.Clean("/"+hex.EncodeToString(hash)))
// 	if errs != nil {
// 		fmt.Println("TestStoreFile failed at reading stored file from stored")
// 		fmt.Println(errs)
// 	} else {
// 		fmt.Println("TestStoreFile was successful")
// 	}

// }

func TestFetchFileRecieve(t *testing.T) {
	RTable := routingTable.NewRoutingTable()
	MBList := &messageBufferList.MessageBufferList{}
	MData := metadata.NewFileMetaData()
	var kadem = NewKademliaObject(RTable, MBList, MData)

	rt := RTable
	for i := 0; i < d7024e.IDLength*8; i++ {
		rt.Buckets[i] = d7024e.NewBucket()
	}
	rt.Me.ID = d7024e.NewKademliaID("0000000000000000000000000000000000000000")

	net := testNetwork{}
	kadem.SetNetworkHandler(&net)
	rpcID1 := d7024e.NewRandomKademliaID()
	rpcID2 := d7024e.NewRandomKademliaID()

	fileHash1 := d7024e.NewRandomKademliaID()
	fileHash2 := d7024e.NewRandomKademliaID()

	checkdata := make(map[string]bool)

	senderIP := "localhost:8001"
	otherHostIP := "localhost:8002"

	checkdata["sender"] = false
	checkdata["ip"] = false

	var wg sync.WaitGroup

	wg.Add(2)

	var cList []testNetworkFetchControl
	cList = append(cList, testNetworkFetchControl{
		func(addr string, filePath string) {
			if addr == senderIP+"/storage/"+fileHash1.String() {
				checkdata["sender"] = true
			} else if addr == otherHostIP+"/storage/"+fileHash1.String() {
				checkdata["ip"] = true
			} else {
				panic("IP NOT RECOGNIZED")
			}
			wg.Done()
		},
	})
	cList = append(cList, testNetworkFetchControl{
		func(addr string, filePath string) {
			if addr == senderIP+"/storage/"+fileHash2.String() {
				checkdata["sender"] = true
			} else if addr == otherHostIP+"/storage/"+fileHash2.String() {
				checkdata["ip"] = true
			} else {
				panic("IP NOT RECOGNIZED")
			}
			wg.Done()
		},
	})

	net.CheckListFetch = cList

	senderStoreData := rpc.StoreFile{*fileHash1, "sender"}
	senderStoreByteData, err := json.Marshal(senderStoreData)
	if err != nil {
		panic("MARSHALING ERROR")
	}
	senderStoreMessage := rpc.Message{rpc.STORE, *rpcID1, *d7024e.NewKademliaID("0000000000000000000000000000000000000001"), senderStoreByteData}
	senderStoreByteMessage, err := json.Marshal(senderStoreMessage)
	if err != nil {
		panic("MARSHALING ERROR")
	}
	otherHostStoreData := rpc.StoreFile{*fileHash2, otherHostIP}
	otherHostStoreByteData, err := json.Marshal(otherHostStoreData)
	if err != nil {
		panic("MARSHALING ERROR")
	}
	otherHostStoreMessage := rpc.Message{rpc.STORE, *rpcID2, *d7024e.NewKademliaID("0000000000000000000000000000000000000002"), otherHostStoreByteData}
	otherHostStoreByteMessage, err := json.Marshal(otherHostStoreMessage)
	if err != nil {
		panic("MARSHALING ERROR")
	}
	kadem.HandleIncomingRPC(senderStoreByteMessage, senderIP)
	kadem.HandleIncomingRPC(otherHostStoreByteMessage, otherHostIP)

	wg.Wait()

	if checkdata["sender"] == true && checkdata["ip"] == true {
		println("TestFetchFileRecieve successful")
	} else {
		panic("TestFetchFileRecieve failed")
	}

}

func TestAddContact(t *testing.T) {
	RTable := routingTable.NewRoutingTable()
	MBList := &messageBufferList.MessageBufferList{}
	MData := metadata.NewFileMetaData()
	var kadem = NewKademliaObject(RTable, MBList, MData)

	rt := RTable
	for i := 0; i < d7024e.IDLength*8; i++ {
		rt.Buckets[i] = d7024e.NewBucket()
	}
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

			kadem.HandleIncomingRPC(d, addr)
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

			mbList := MBList
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
	RTable := routingTable.NewRoutingTable()
	MBList := &messageBufferList.MessageBufferList{}
	MData := metadata.NewFileMetaData()
	var kadem = NewKademliaObject(RTable, MBList, MData)

	fmt.Println("ping test")
	rt := RTable
	for i := 0; i < d7024e.IDLength*8; i++ {
		rt.Buckets[i] = d7024e.NewBucket()
	}
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

func TestReturnContacts(t *testing.T) {
	RTable := routingTable.NewRoutingTable()
	MBList := &messageBufferList.MessageBufferList{}
	MData := metadata.NewFileMetaData()
	var kadem = NewKademliaObject(RTable, MBList, MData)

	rt := RTable
	sender := d7024e.NewContact(d7024e.NewKademliaID("FFFFFFFFF0000000000FF0000000000000000001"), "localhost:9999")
	rpcId := d7024e.NewRandomKademliaID()

	for i := 0; i < d7024e.IDLength*8; i++ {
		rt.Buckets[i] = d7024e.NewBucket()
	}

	startingContacts := []d7024e.Contact{
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

	for _, c := range startingContacts {
		rt.AddContact(c)
	}

	var cList []testNetworkControl
	cList = append(cList, testNetworkControl{
		func(sentMessage rpc.Message, addr string) {
			assertEqual(t, rpc.CLOSEST_NODES, sentMessage.RpcType)
			assertEqual(t, addr, sender.Address)
			var closestNodes rpc.ClosestNodes
			json.Unmarshal(sentMessage.RpcData, &closestNodes)
			for i, c := range closestNodes.Closest {
				assertEqual(t, c.ID.String(), startingContacts[i].ID.String())
			}
		},
	})
	net := testNetwork{}
	net.CheckList = cList
	kadem.SetNetworkHandler(&net)

	rpcData := rpc.FindNode{*d7024e.NewKademliaID("FFFFFFFFF0000000000000000000000000000000")}
	kadem.handleFindNode(*rpcId, rpcData, sender.Address)

}

func TestHandleFindValue(t *testing.T) {

	RTable := routingTable.NewRoutingTable()
	MBList := &messageBufferList.MessageBufferList{}
	MData := metadata.NewFileMetaData()
	var kadem = NewKademliaObject(RTable, MBList, MData)

	fmt.Println("findvalue 1")
	target := d7024e.NewContact(d7024e.NewKademliaID("FFFFFFFFF0000000000000000000000000000000"), "localhost:8000")
	lookUpTestInitalSetup(target, t, kadem, RTable)

	rpcID := d7024e.NewRandomKademliaID()
	senderID := d7024e.NewKademliaID("FFFFFFFFF000000000000000000000FF00000007")
	rpcData := rpc.FindNode{*d7024e.NewKademliaID("FFFFFFFFF0000000000000000000000000000000")}
	rpcM, _ := rpc.Marshal(rpc.FIND_VALUE, *rpcID, *senderID, rpcData)

	metaData := MData
	metaData.AddFile("N/A", "fffffffff0000000000000000000000000000000", false, time.Hour)

	cList := []testNetworkControl{
		testNetworkControl{
			func(sentMessage rpc.Message, addr string) {
				if sentMessage.RpcType == rpc.HAS_VALUE {
					if !sentMessage.RpcId.Equals(rpcID) {
						fmt.Println("Invalid rpc id")
						t.Fail()
					}
				} else {
					fmt.Println("Invalid rpc type")
					t.Fail()
				}
			},
		},
	}

	net := testNetwork{}
	net.CheckList = cList
	kadem.SetNetworkHandler(&net)
	fmt.Println("Started")
	kadem.HandleIncomingRPC(rpcM, "localhost:9999")
	fmt.Println("Came back")
	time.Sleep(2 * time.Second)
}

func TestHandleFindValueNotFound(t *testing.T) {
	RTable := routingTable.NewRoutingTable()
	MBList := &messageBufferList.MessageBufferList{}
	MData := metadata.NewFileMetaData()
	var kadem = NewKademliaObject(RTable, MBList, MData)

	fmt.Println("Started find value")
	target := d7024e.NewContact(d7024e.NewKademliaID("FFFFFFFFF0000000000000000000000000000000"), "localhost:8000")
	lookUpTestInitalSetup(target, t, kadem, RTable)

	expectedResult := []d7024e.Contact{
		d7024e.NewContact(d7024e.NewKademliaID("FFFFFFFFF0000000000000000000000000000001"), "localhost:8001"),
		d7024e.NewContact(d7024e.NewKademliaID("FFFFFFFFF0000000000000000000000000000002"), "localhost:8002"),
		d7024e.NewContact(d7024e.NewKademliaID("FFFFFFFFF0000000000000000000000000000003"), "localhost:8003"),
		d7024e.NewContact(d7024e.NewKademliaID("FFFFFFFFF0000000000000000000000000000004"), "localhost:8004"),
		d7024e.NewContact(d7024e.NewKademliaID("FFFFFFFFF0000000000000000000000000000005"), "localhost:8005"),
		d7024e.NewContact(d7024e.NewKademliaID("FFFFFFFFF0000000000000000000000000000006"), "localhost:8006"),
		d7024e.NewContact(d7024e.NewKademliaID("0FFFFFFFF0000000000000000000000000000015"), "localhost:8021"),
		d7024e.NewContact(d7024e.NewKademliaID("0FFFFFFFF0000000000000000000000000000016"), "localhost:8022"),
		d7024e.NewContact(d7024e.NewKademliaID("0FFFFFFFF0000000000000000000000000000017"), "localhost:8023"),
		d7024e.NewContact(d7024e.NewKademliaID("0FFFFFFFF0000000000000000000000000000018"), "localhost:8024"),
		d7024e.NewContact(d7024e.NewKademliaID("0FFFFFFFF0000000000000000000000000000019"), "localhost:8025"),
		d7024e.NewContact(d7024e.NewKademliaID("0FFFFFFFF000000000000000000000000000001A"), "localhost:8026"),
		d7024e.NewContact(d7024e.NewKademliaID("0FFFFFFFF000000000000000000000000000001B"), "localhost:8027"),
		d7024e.NewContact(d7024e.NewKademliaID("0FFFFFFFF000000000000000000000000000001C"), "localhost:8028"),
		d7024e.NewContact(d7024e.NewKademliaID("0FFFFFFFF000000000000000000000000000001D"), "localhost:8029"),
	}

	rpcID := d7024e.NewRandomKademliaID()
	senderID := d7024e.NewKademliaID("FFFFFFFFF0000000000000000000000000000001")
	rpcData := rpc.FindNode{*d7024e.NewKademliaID("FFFFFFFFF00000000000000000ff000000000000")}
	rpcM, _ := rpc.Marshal(rpc.FIND_VALUE, *rpcID, *senderID, rpcData)

	cList := []testNetworkControl{
		testNetworkControl{
			func(sentMessage rpc.Message, addr string) {
				assertEqual(t, addr, "localhost:9999")
				if sentMessage.RpcType == rpc.CLOSEST_NODES {
					if !sentMessage.RpcId.Equals(rpcID) {
						t.Fail()
						fmt.Println("Unexpected rpc id")
					}

					var reply rpc.ClosestNodes
					json.Unmarshal(sentMessage.RpcData, &reply)

					if len(reply.Closest) == len(expectedResult) {
						for i, c := range reply.Closest {
							assertEqual(t, c.ID.String(), expectedResult[i].ID.String())
						}
					} else {
						fmt.Println("Unexpected amount of contacts")
						t.Fail()
					}
				} else {
					fmt.Println("Unexpected rpc type")
					t.Fail()
				}
			},
		},
	}

	net := testNetwork{}
	net.CheckList = cList
	kadem.SetNetworkHandler(&net)
	kadem.HandleIncomingRPC(rpcM, "localhost:9999")
	time.Sleep(3 * time.Second)
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
