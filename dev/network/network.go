package network

import (
	"Kademlia_D7024E/dev/d7024e"
	"Kademlia_D7024E/dev/messageBufferList"
	"Kademlia_D7024E/dev/routingTable"
	"Kademlia_D7024E/dev/rpc"
	"encoding/json"
	"fmt"
	"net"
	"strconv"
	"sync"
)

const MAX_PACKET_SIZE int = 1 //TODO Calculate real max packet size.

type network struct {
	port int
	ip   string
}

var instance *network
var once sync.Once

func GetInstance() *network {
	once.Do(func() {
		instance = &network{}
	})
	return instance
}

func SetPort(port int) {
	GetInstance().port = port
}

func SetIp(ip string) {
	GetInstance().ip = ip
}

func Listen(ip string, port int) {
	serverAddr, addrErr := net.ResolveUDPAddr("udp", ip+":"+strconv.Itoa(port))
	if addrErr != nil {
		return
	}
	for {
		conn, listenErr := net.ListenUDP("udp", serverAddr)
		if listenErr != nil {
			continue
		}
		go HandleConnection(conn)
	}

}

func HandleConnection(conn *net.UDPConn) {
	var data []byte
	var message rpc.Message = rpc.Message{}
	_, _, err := conn.ReadFrom(data)

	if err != nil {
		//TODO handle error
	}

	unmarshaling_err := json.Unmarshal(data, &message)
	if unmarshaling_err != nil {

	}

	switch message.RpcType {

	case rpc.FIND_NODE:
		var find_node rpc.FindNode
		json.Unmarshal(message.RpcData, find_node)
		HandleFindNode(message.RpcId, find_node, conn)

	case rpc.PING:
		HandlePing(message.RpcId, conn)

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

func HandleFindNode(rpc_id d7024e.KademliaID, find_node rpc.FindNode, conn *net.UDPConn) {
	rt := routingTable.GetInstance()
	closest_nodes := rpc.ClosestNodes{rt.FindClosestContacts(&find_node.NodeId, 20)}
	response := rpc.Marshal(rpc.CLOSEST_NODES, rpc_id, closest_nodes)

	conn.Write(response)
	conn.Close()
}

func HandlePing(rpc_id d7024e.KademliaID, conn *net.UDPConn) {
	response, err := json.Marshal(rpc.Message{RpcType: rpc.PONG, RpcId: rpc_id, RpcData: []byte{byte(0)}})
	if err != nil {
		fmt.Println(err)
		return
	}
	conn.Write(response)
	conn.Close()
}

func (network *network) SendPingMessage(contact *d7024e.Contact, rpc_id *d7024e.KademliaID) {
	data, m_err := json.Marshal(rpc.Message{RpcType: rpc.PING, RpcId: *rpc_id, RpcData: []byte{byte(0)}})

	if m_err != nil {
		fmt.Println(m_err)
		return
	}

	laddr, l_err := net.ResolveUDPAddr("udp", network.ip+":"+strconv.Itoa(network.port))
	raddr, r_err := net.ResolveUDPAddr("udp", contact.Address)

	if l_err != nil || r_err != nil {
		if l_err != nil {
			fmt.Println(l_err)
		} else if l_err != nil {
			fmt.Println(r_err)
		}
		return
	}

	conn, err := net.DialUDP("udp", laddr, raddr)
	if err != nil {
		fmt.Println(err)
		return
	}
	conn.Write(data)
}

func (network *network) SendFindContactMessage(contact *d7024e.Contact, toFind *d7024e.KademliaID, rpc_id *d7024e.KademliaID) {

	data := rpc.Marshal(rpc.FIND_NODE, *rpc_id, rpc.FindNode{*toFind})

	laddr, l_err := net.ResolveUDPAddr("udp", network.ip+":"+strconv.Itoa(network.port))
	raddr, r_err := net.ResolveUDPAddr("udp", contact.Address)

	if l_err != nil || r_err != nil {
		if l_err != nil {
			fmt.Println(l_err)
		} else if l_err != nil {
			fmt.Println(r_err)
		}
		return
	}

	conn, err := net.DialUDP("udp", laddr, raddr)
	if err != nil {
		fmt.Println(err)
		return
	}
	conn.Write(data)
}

func (network *network) SendFindDataMessage(contact *d7024e.Contact, toFind *d7024e.KademliaID, rpc_id *d7024e.KademliaID) {
	data := rpc.Marshal(rpc.FIND_VALUE, *rpc_id, rpc.FindNode{*toFind})

	laddr, l_err := net.ResolveUDPAddr("udp", network.ip+":"+strconv.Itoa(network.port))
	raddr, r_err := net.ResolveUDPAddr("udp", contact.Address)

	if l_err != nil || r_err != nil {
		if l_err != nil {
			fmt.Println(l_err)
		} else if l_err != nil {
			fmt.Println(r_err)
		}
		return
	}

	conn, err := net.DialUDP("udp", laddr, raddr)
	if err != nil {
		fmt.Println(err)
		return
	}
	conn.Write(data)
}

func (network *network) SendStoreMessage(data []byte) {
	// TODO
}
