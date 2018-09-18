package rpc

import (
	"Kademlia_D7024E/dev/d7024e"
	"encoding/json"
)

const FIND_NODE = "find_node"
const CLOSEST_NODES = "closest_nodes"
const FIND_VALUE = "find_value"
const STORE = "store"
const PING = "ping"
const PONG = "pong"

type Message struct {
	RpcType string
	RpcId   d7024e.KademliaID
	RpcData []byte
}

type FindNode struct {
	NodeId d7024e.KademliaID
}

type ClosestNodes struct {
	Closest []d7024e.Contact
}

func Marshal(rpc_type string, rpc_id d7024e.KademliaID, rpc_data interface{}) ([]byte, error) {
	m_rpc_data, data_err := json.Marshal(rpc_data)
	if data_err != nil {
		return []byte{}, data_err
	}

	new_message := Message{rpc_type, rpc_id, m_rpc_data}
	m_message, message_err := json.Marshal(new_message)
	if message_err != nil {
		return []byte{}, message_err
	}

	return m_message, nil
}
