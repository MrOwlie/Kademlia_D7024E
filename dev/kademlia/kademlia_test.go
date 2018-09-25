package kademlia

import (
	"testing"

	"../rpc"
)

type testNetworkControl struct {
	SentMsg   rpc.Message
	ReturnMsg rpc.Message
}

type testNetwork struct {
	CheckList []testNetworkControl
}

func (net *testNetwork) Listen() {

}

func (net *testNetwork) SendMessage(addr string, data *[]byte) {

}

func TestFindNode(t *testing.T) {

}

func TestJoin(t *testing.T) {
	/*	kadem := GetInstance()
		var cList []testNetworkControl = {
			testNetworkControl{
				rpc.Message{rpc.FIND_NODE, }
			}
		}
		net := testNetwork{

		}
		kadem.SetNetworkHandler()*/
}
