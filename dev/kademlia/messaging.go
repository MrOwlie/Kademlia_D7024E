package kademlia

import (
	"encoding/json"
	"fmt"
	"os"
	"encoding/hex"

	"../d7024e"
	"../messageBufferList"
	"../network"
	"../routingTable"
	"../rpc"
)

var storagePath string = "What ever the storage path is" //TODO fix this

func (kademlia *kademlia) HandleIncomingRPC(data []byte, addr string) {
	var message rpc.Message = rpc.Message{}
	unmarshaling_err := json.Unmarshal(data, &message)
	if unmarshaling_err != nil {
		fmt.Println(unmarshaling_err)
	}

	//Update contact
	routingTable.GetInstance().AddContact(d7024e.NewContact(&message.SenderId, addr))
	switch message.RpcType {

	case rpc.FIND_NODE:
		var find_node rpc.FindNode
		json.Unmarshal(message.RpcData, &find_node)
		kademlia.handleFindNode(message.RpcId, find_node, addr)

	case rpc.PING:
		kademlia.handlePing(message.RpcId, addr)

	case rpc.FIND_VALUE:
		var find_node rpc.FindNode
		json.Unmarshal(message.RpcData, &find_node)
		kademlia.handleFindValue(message.RpcId, find_node, addr)

	default:
		if message.RpcType == rpc.CLOSEST_NODES || message.RpcType == rpc.PONG {
			buffer_list := messageBufferList.GetInstance()
			m_buffer, hasId := buffer_list.GetMessageBuffer(&message.RpcId)
			if hasId {
				m_buffer.AppendMessage(message)
			} else {
				fmt.Println("Message with rpc id: %v was discarded", message.RpcId)
			}
		} else {
			fmt.Println("Invalid message")
		}
	}
}

func (kademlia *kademlia) handleFindNode(rpc_id d7024e.KademliaID, find_node rpc.FindNode, addr string) {
	rt := routingTable.GetInstance()
	closest_nodes := rpc.ClosestNodes{rt.FindClosestContacts(&find_node.NodeId, 20)}
	response, err := rpc.Marshal(rpc.CLOSEST_NODES, rpc_id, *rt.Me.ID, closest_nodes)
	fmt.Println("sending contacts: ", closest_nodes)

	if err != nil {
		fmt.Println(err)
	}

	network.GetInstance().SendMessage(addr, response)
}

func (kademlia *kademlia) handlePing(rpc_id d7024e.KademliaID, addr string) {
	rt := routingTable.GetInstance()
	response, err := json.Marshal(rpc.Message{RpcType: rpc.PONG, RpcId: rpc_id, SenderId: *rt.Me.ID, RpcData: []byte{byte(0)}})
	if err != nil {
		fmt.Println(err)
		return
	}

	network.GetInstance().SendMessage(addr, response)
}

func (kademlia *kademlia) handleFindValue(rpc_id d7024e.KademliaID, find_node rpc.FindNode, addr string){
	rt := routingTable.GetInstance()
	fileName := hex.EncodeToString(find_node.NodeId)
	var response []byte

	if _, err := os.Stat(storagePath+"/"+fileName); err != nil {
		if !os.IsNotExist(err) {
			fmt.Println(err)
		} else {
			closest_nodes := rpc.ClosestNodes{rt.FindClosestContacts(&find_node.NodeId, 20)}
			response, err = rpc.Marshal(rpc.CLOSEST_NODES, rpc_id, *rt.Me.ID, closest_nodes)
			
			if err != nil {
				fmt.Println(err)
			}
		}
	} else {
		response, err = json.Marshal(rpc.Message{rpc.HAS_VALUE, rpc_id, *rt.Me.ID, []byte{byte(0)}})

		if err != nil {
			fmt.Println(err)
		}
	}

	network.GetInstance().SendMessage(addr, response)
}

func (kademlia *kademlia) sendPingMessage(contact *d7024e.Contact, rpc_id *d7024e.KademliaID) {
	rt := routingTable.GetInstance()
	data, m_err := json.Marshal(rpc.Message{RpcType: rpc.PING, RpcId: *rpc_id, SenderId: *rt.Me.ID, RpcData: []byte{byte(0)}})

	if m_err != nil {
		fmt.Println(m_err)
		return
	}

	/*addr, addr_err := net.ResolveUDPAddr("udp", contact.Address)
	if addr_err != nil {
		fmt.Println(addr_err)
	}*/

	network.GetInstance().SendMessage(contact.Address, &data)
}

func (kademlia *kademlia) sendFindContactMessage(contact *d7024e.Contact, toFind *d7024e.KademliaID, rpc_id *d7024e.KademliaID) {
	rt := routingTable.GetInstance()
	data, err := rpc.Marshal(rpc.FIND_NODE, *rpc_id, *rt.Me.ID, rpc.FindNode{*toFind})
	if err != nil {
		fmt.Println(err)
	}

	/*addr, addr_err := net.ResolveUDPAddr("udp", contact.Address)
	if addr_err != nil {
		fmt.Println(addr_err)
	}*/

	network.GetInstance().SendMessage(contact.Address, &data)
}

func (kademlia *kademlia) sendFindDataMessage(contact *d7024e.Contact, toFind *d7024e.KademliaID, rpc_id *d7024e.KademliaID) {
	rt := routingTable.GetInstance()
	data, err := rpc.Marshal(rpc.FIND_VALUE, *rpc_id, *rt.Me.ID, rpc.FindNode{*toFind})
	if err != nil {
		fmt.Println(err)
	}

	/*addr, addr_err := net.ResolveUDPAddr("udp", contact.Address)
	if addr_err != nil {
		fmt.Println(addr_err)
	}*/

	network.GetInstance().SendMessage(contact.Address, &data)
}

func (kademlia *kademlia) sendStoreMessage(contact *d7024e.Contact, rpc_id *d7024e.KademliaID, fileHash *d7024e.KademliaID) {
	rt := routingTable.GetInstance()
	data, err := rpc.Marshal(rpc.STORE, *rpc_id, *rt.Me.ID, rpc.StoreFile{*fileHash})
	if err != nil {
		fmt.Println(err)
	}

	network.GetInstance().SendMessage(contact.Address, data) 
}
