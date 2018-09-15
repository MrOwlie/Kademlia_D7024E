package rpc

import (
	"encoding/json"
	"./d7024e"
	)

const FIND_NODE = "find_node"
const CLOSEST_NODES = "closest_nodes"
const STORE = "store"
const PING = "ping"
const PONG = "pong"

type message struct {
	rpc_type string
	rpc_id int
	rpc_data json.RawMessage
}

type find_node{
	node_id int
}

type closest_nodes struct{
	closest_k d7024e.Contact[]
}

func Marshal(rpc_type string, rpc_id int, rpc_data interface{}) byte[]{
	m_rpc_data, data_err := json.Marshal(rpc_data)
	if(data_err != nil){
		// TODO: error handeling
	}

	new_message := message{rpc_type, rpc_id, rpc_data}
	m_message, message_err := json.Marshal(new_message)
	if(message_err != nil){
		// TODO: error handeling
	}

	return m_message;
}
