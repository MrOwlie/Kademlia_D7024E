package kademlia

import (
	"../rpc"
)

type DummyNetwork struct{
	m map[string]rpc.Message
}