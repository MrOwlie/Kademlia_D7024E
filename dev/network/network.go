package network

import (
	"sync"
	"../d7024e"
	"net"
	"sync"
	"rpc"
)


type network struct {
	port int
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

func Listen(ip string, port int) {
	serverAddr, addrErr = net.ResolveUDPAddr("udp", ip+":"+port)
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

func HandleConnection(*UDPConn conn){

}

func (network *network) SendPingMessage(contact *d7024e.Contact) {
	// TODO
}

func (network *network) SendFindContactMessage(contact *d7024e.Contact, toFind *d7024e.KademliaID) {
	// TODO
}

func (network *network) SendFindDataMessage(hash string) {
	// TODO
}

func (network *network) SendStoreMessage(data []byte) {
	// TODO
}
