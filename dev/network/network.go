package network

import (
	"sync"
	"../d7024e"
	"../routingTable"
	"net"
	"sync"
	"strconv"
	"encoding/json"
	"../messageBufferList"
	"fmt"
)

const MAX_PACKET_SIZE int = 1 //TODO Calculate real max packet size.

type network struct {
	port int
	ip string
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

func SetIp(ip string){
	GetInstance().ip = ip
}

func Listen(ip string, port int) {
	serverAddr, addrErr := net.ResolveUDPAddr("udp", ip+":"+strconv.Itoa(port))
	if addrErr != nil {
		return
	}
	for{
		conn, listenErr = net.ListenUDP("udp", serverAddr)
		if listenErr != nil {
			continue
		}
		go HandleConnection(conn)
	}

}

func HandleConnection(*net.UDPConn conn){
	data [MAX_PACKET_SIZE]byte
	var message Message
	_, _, err := conn.ReadFrom(data)

	if err != nil{
		//TODO handle error
	}

	unmarshaling_err := json.Unmarshal(data, &message)
	if(unmarshaling_err != nil){

	}

	switch message.rpc_type{

		case FIND_NODE:
			var find_node FindNode
			json.Unmarshal(message.rpc_data, find_node)
			HandleFindNode(message.rpc_id, find_node, conn)

		case PING:
			HandlePing(message.rpc_id, conn)

		default:
			if message.rpc_type == CLOSEST_NODES || message.rpc_type == PONG {
				m_buffer, hasId = messageBufferList.GetInstance().GetMessageBuffer(message.rpc_id);
				if(hasId){
					m_buffer.appendMessage(message);
				}
			}
	}

}

func HandleFindNode(rpc_id d7024e.KademliaID, find_node FindNode, conn *net.UDPConn){
	rt := routingTable.GetInstance()
	closest_nodes := ClosestNodes{rt.FindClosestContacts(find_node.NodeId, 20)}
	response := Marshal(CLOSEST_NODES, rpc_id, closest_nodes)

	conn.Write(response)
	conn.Close()
}

func HandlePing(rpc_id d7024e.KademliaID, conn *net.UDPConn){
	var response []byte
	json.Marshal(Message(rpc_type: PONG, rpc_id: rpc_id, rpc_data: byte(0)), response)
	conn.Write(response)
	conn.Close()
}


func (network *network) SendPingMessage(contact *d7024e.Contact, rpc_id *d7024e.KademliaID) {
	data, m_err := json.Marshal(Message(rpc_type: PING, rpc_id: rpc_id, rpc_data: byte(0)))
	
	if m_err != nil{
		fmt.Println(m_err)
		return
	}

	laddr := net.ResolveUDPAddr("udp", network.ip+":"+strconv.Itoa(network.port))
	raddr := net.ResolveUDPAddr("udp", contact.Address)

	conn, err := net.DialUDP("udp", laddr, raddr)
	if(err != nil){
		fmt.Println(err)
		return
	}
	conn.Write(data)
}

func (network *network) SendFindContactMessage(contact *d7024e.Contact, toFind *d7024e.KademliaID, rpc_id *d7024e.KademliaID) {

	data = Marshal(FIND_NODE, rpc_id, FindNode{toFind})
	
	laddr := net.ResolveUDPAddr("udp", network.ip+":"+strconv.Itoa(network.port))
	raddr := net.ResolveUDPAddr("udp", contact.Address)

	conn, err := net.DialUDP("udp", laddr, raddr)
	if(err != nil){
		fmt.Println(err)
		return
	}
	conn.Write(data)
}

func (network *network) SendFindDataMessage(contact *d7024e.Contact, toFind *d7024e.KademliaID, rpc_id *d7024e.KademliaID) {
	data = Marshal(FIND_VALUE, rpc_id, FindNode{toFind})
	
	laddr := net.ResolveUDPAddr("udp", network.ip+":"+strconv.Itoa(network.port))
	raddr := net.ResolveUDPAddr("udp", contact.Address)

	conn, err := net.DialUDP("udp", laddr, raddr)
	if(err != nil){
		fmt.Println(err)
		return
	}
	conn.Write(data)
}

func (network *network) SendStoreMessage(data []byte) {
	// TODO
}
