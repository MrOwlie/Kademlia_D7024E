package network

import (
	"encoding/json"
	"fmt"
	"net"
	"strconv"
	"sync"

	"../d7024e"
	"../messageBufferList"
	"../routingTable"
	"../rpc"
)

const MAX_PACKET_SIZE int = 512 //TODO Calculate actual max packet size.

type Handler interface {
	HandleFindNode()
}

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

	conn, listenErr := net.ListenUDP("udp", serverAddr)
	if listenErr != nil {
		return
	}

	fmt.Println("Listening to UDP traffic on " + ip + ":" + strconv.Itoa(port))
	for {
		var data [MAX_PACKET_SIZE]byte
		n, addr, err := conn.ReadFromUDP(data[0:])
		if err != nil {
			fmt.Println(err)
			continue
		}
		go HandleConnection(data[0:n], addr)
	}

}

func HandleConnection(data []byte, addr *net.UDPAddr) {
	var message rpc.Message = rpc.Message{}
	unmarshaling_err := json.Unmarshal(data, &message)
	if unmarshaling_err != nil {
		fmt.Println(unmarshaling_err)
	}

	switch message.RpcType {

	case rpc.FIND_NODE:
		var find_node rpc.FindNode
		json.Unmarshal(message.RpcData, &find_node)
		HandleFindNode(message.RpcId, find_node, addr)

	case rpc.PING:
		HandlePing(message.RpcId, addr)

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

func HandleFindNode(rpc_id d7024e.KademliaID, find_node rpc.FindNode, addr *net.UDPAddr) {
	rt := routingTable.GetInstance()
	closest_nodes := rpc.ClosestNodes{rt.FindClosestContacts(&find_node.NodeId, 20)}
	response, err := rpc.Marshal(rpc.CLOSEST_NODES, rpc_id, closest_nodes)

	if err != nil {
		fmt.Println(err)
	}

	GetInstance().sendMessage(addr, response)
}

func HandlePing(rpc_id d7024e.KademliaID, addr *net.UDPAddr) {
	response, err := json.Marshal(rpc.Message{RpcType: rpc.PONG, RpcId: rpc_id, RpcData: []byte{byte(0)}})
	if err != nil {
		fmt.Println(err)
		return
	}

	GetInstance().sendMessage(addr, response)
}

func (network *network) SendPingMessage(contact *d7024e.Contact, rpc_id *d7024e.KademliaID) {
	data, m_err := json.Marshal(rpc.Message{RpcType: rpc.PING, RpcId: *rpc_id, RpcData: []byte{byte(0)}})

	if m_err != nil {
		fmt.Println(m_err)
		return
	}

	addr, addr_err := net.ResolveUDPAddr("udp", contact.Address)
	if addr_err != nil {
		fmt.Println(addr_err)
	}

	network.sendMessage(addr, data)
}

func (network *network) SendFindContactMessage(contact *d7024e.Contact, toFind *d7024e.KademliaID, rpc_id *d7024e.KademliaID) {
	data, err := rpc.Marshal(rpc.FIND_NODE, *rpc_id, rpc.FindNode{*toFind})
	if err != nil {
		fmt.Println(err)
	}

	addr, addr_err := net.ResolveUDPAddr("udp", contact.Address)
	if addr_err != nil {
		fmt.Println(addr_err)
	}

	network.sendMessage(addr, data)
}

func (network *network) SendFindDataMessage(contact *d7024e.Contact, toFind *d7024e.KademliaID, rpc_id *d7024e.KademliaID) {
	data, err := rpc.Marshal(rpc.FIND_VALUE, *rpc_id, rpc.FindNode{*toFind})
	if err != nil {
		fmt.Println(err)
	}

	addr, addr_err := net.ResolveUDPAddr("udp", contact.Address)
	if addr_err != nil {
		fmt.Println(addr_err)
	}

	network.sendMessage(addr, data)
}

func (network *network) SendStoreMessage(data []byte) {
	// TODO
}

func (network *network) sendMessage(addr *net.UDPAddr, data []byte) {

	laddr, l_err := net.ResolveUDPAddr("udp", network.ip+":"+strconv.Itoa(network.port))

	if l_err != nil {
		fmt.Println(l_err)
		return
	}

	conn, err := net.DialUDP("udp", laddr, addr)
	if err != nil {
		fmt.Println(err)
		return
	}
	conn.Write(data)
	conn.Close()
}
