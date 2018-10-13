package kademlia

import (
	"encoding/json"
	"fmt"

	"../d7024e"
	"../rpc"
)

//var storagePath string = "What ever the storage path is" //TODO fix this

func (kademlia *Kademlia) HandleIncomingRPC(data []byte, addr string) {
	//fmt.Println("msg ", string(data))
	var message rpc.Message = rpc.Message{}
	unmarshaling_err := json.Unmarshal(data, &message)
	if unmarshaling_err != nil {
		fmt.Println(unmarshaling_err)
	}

	//Update contact
	con := d7024e.NewContact(&message.SenderId, addr)
	kademlia.addContact(&con)
	fmt.Println("adding contact ", addr, " ", &message.SenderId)
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

	case rpc.STORE:
		var store_file rpc.StoreFile
		json.Unmarshal(message.RpcData, &store_file)
		kademlia.handleStore(&store_file, addr)

	/*case rpc.TIME_OUT:
	kademlia.HandleTimeout(message.RpcId)*/

	default:
		if message.RpcType == rpc.CLOSEST_NODES || message.RpcType == rpc.PONG || message.RpcType == rpc.HAS_VALUE {
			buffer_list := kademlia.MBList
			m_buffer, hasId := buffer_list.GetMessageBuffer(&message.RpcId)
			if hasId {
				m_buffer.AppendMessage(&message)
			} else {
				fmt.Printf("Message with rpc id: %v was discarded", message.RpcId)
			}
		} else {
			fmt.Println("Invalid message")
		}
	}
}

func (kademlia *Kademlia) handleFindNode(rpc_id d7024e.KademliaID, find_node rpc.FindNode, addr string) {
	rt := kademlia.routingTable
	closest_nodes := rpc.ClosestNodes{rt.FindClosestContacts(&find_node.NodeId, 20)}
	response, err := rpc.Marshal(rpc.CLOSEST_NODES, rpc_id, *rt.Me.ID, closest_nodes)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("sending closest contacts ", len(closest_nodes.Closest))

	kademlia.network.SendMessage(addr, &response)
}

func (kademlia *Kademlia) handlePing(rpc_id d7024e.KademliaID, addr string) {
	rt := kademlia.routingTable
	response, err := json.Marshal(rpc.Message{RpcType: rpc.PONG, RpcId: rpc_id, SenderId: *rt.Me.ID, RpcData: []byte{byte(0)}})
	if err != nil {
		fmt.Println(err)
		return
	}

	kademlia.network.SendMessage(addr, &response)
}

func (kademlia *Kademlia) handleFindValue(rpc_id d7024e.KademliaID, find_node rpc.FindNode, addr string) {
	rt := kademlia.routingTable
	metadata := kademlia.MetaData

	hash := find_node.NodeId.String()
	var response []byte
	var err error

	fmt.Println(hash)
	if metadata.HasFile(hash) {
		fmt.Println("found file")
		response, err = json.Marshal(rpc.Message{rpc.HAS_VALUE, rpc_id, *rt.Me.ID, []byte{byte(0)}})

		if err != nil {
			fmt.Println(err)
		}
	} else {
		fmt.Println("did not find file")
		closest_nodes := rpc.ClosestNodes{rt.FindClosestContacts(&find_node.NodeId, 20)}
		response, err = rpc.Marshal(rpc.CLOSEST_NODES, rpc_id, *rt.Me.ID, closest_nodes)

		if err != nil {
			fmt.Println(err)
		}
	}

	kademlia.network.SendMessage(addr, &response)
}

func (kademlia *Kademlia) handleStore(store_file *rpc.StoreFile, addr string) {
	metadata := kademlia.MetaData
	hash := store_file.FileHash.String()

	if metadata.HasFile(hash) {
		metadata.RefreshFile(hash)
	} else {
		var hostURL string
		filePath := storagePath + "/" + hash

		if store_file.Host == rpc.SENDER {
			hostURL = addr
		} else {
			hostURL = store_file.Host
		}

		hostURL += "/storage/" + hash
		err := kademlia.network.FetchFile(hostURL, filePath)
		if err == nil {
			metadata.AddFile(filePath, hash, false, kademlia.calcTimeToLive(&store_file.FileHash))
		} else {
			fmt.Println(err, " (", hostURL, ")")
		}
	}
}

func (kademlia *Kademlia) sendPingMessage(contact *d7024e.Contact, rpc_id *d7024e.KademliaID) {
	rt := kademlia.routingTable
	data, m_err := json.Marshal(rpc.Message{RpcType: rpc.PING, RpcId: *rpc_id, SenderId: *rt.Me.ID, RpcData: []byte{byte(0)}})

	if m_err != nil {
		fmt.Println(m_err)
		return
	}

	/*addr, addr_err := net.ResolveUDPAddr("udp", contact.Address)
	if addr_err != nil {
		fmt.Println(addr_err)
	}*/

	kademlia.network.SendMessage(contact.Address, &data)
}

func (kademlia *Kademlia) sendFindContactMessage(contact *d7024e.Contact, toFind *d7024e.KademliaID, rpc_id *d7024e.KademliaID) {
	rt := kademlia.routingTable
	data, err := rpc.Marshal(rpc.FIND_NODE, *rpc_id, *rt.Me.ID, rpc.FindNode{*toFind})
	if err != nil {
		fmt.Println(err)
	}

	/*addr, addr_err := net.ResolveUDPAddr("udp", contact.Address)
	if addr_err != nil {
		fmt.Println(addr_err)
	}*/

	kademlia.network.SendMessage(contact.Address, &data)
}

func (kademlia *Kademlia) sendFindDataMessage(contact *d7024e.Contact, toFind *d7024e.KademliaID, rpc_id *d7024e.KademliaID) {
	rt := kademlia.routingTable
	data, err := rpc.Marshal(rpc.FIND_VALUE, *rpc_id, *rt.Me.ID, rpc.FindNode{*toFind})
	if err != nil {
		fmt.Println(err)
	}

	/*addr, addr_err := net.ResolveUDPAddr("udp", contact.Address)
	if addr_err != nil {
		fmt.Println(addr_err)
	}*/

	kademlia.network.SendMessage(contact.Address, &data)
}

func (kademlia *Kademlia) sendStoreMessage(contact *d7024e.Contact, rpc_id *d7024e.KademliaID, fileHash *d7024e.KademliaID, host string) {
	rt := kademlia.routingTable
	data, err := rpc.Marshal(rpc.STORE, *rpc_id, *rt.Me.ID, rpc.StoreFile{*fileHash, host})
	if err != nil {
		fmt.Println(err)
	}

	kademlia.network.SendMessage(contact.Address, &data)
}
