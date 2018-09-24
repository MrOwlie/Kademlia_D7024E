package network

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strconv"
	"sync"
)

const MAX_PACKET_SIZE int = 5120 //TODO Calculate actual max packet size.
const storagePath string = "/kademlia/storage"

type Handler interface {
	HandleIncomingRPC([]byte, string)
}

type network struct {
	port      int
	ip        string
	conn      *net.UDPConn
	msgHandle Handler
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

func SetHandler(h Handler) {
	GetInstance().msgHandle = h
}

func (network *network) Listen() {
	serverAddr, addrErr := net.ResolveUDPAddr("udp", ":"+strconv.Itoa(network.port))
	if addrErr != nil {
		return
	}

	conn, listenErr := net.ListenUDP("udp", serverAddr)
	if listenErr != nil {
		return
	}
	network.conn = conn
	defer conn.Close()

	fmt.Println("Listening to UDP traffic on port " + strconv.Itoa(network.port))
	for {
		var data [MAX_PACKET_SIZE]byte
		n, addr, err := conn.ReadFromUDP(data[0:])
		if err != nil {
			fmt.Println(err)
			continue
		}
		strAddr := addr.IP.String() + ":" + strconv.Itoa(addr.Port)
		go network.msgHandle.HandleIncomingRPC(data[0:n], strAddr)
	}

}

func (network *network) SendMessage(addr string, data *[]byte) {

	laddr, l_err := net.ResolveUDPAddr("udp", addr)

	if l_err != nil {
		fmt.Println(l_err)
		return
	}

	_, err := network.conn.WriteToUDP(*data, laddr)
	if err != nil {
		fmt.Println("dilili", err)
		return
	}
	//conn.Write(data)
	//conn.Close()
}

func ListenFileServer(ip string, port int) {
	http.Handle("/", http.FileServer(http.Dir(storagePath)))
	err := http.ListenAndServe(":"+strconv.Itoa(port), nil)
	if err != nil {
		fmt.Println(err)
	}
}

func FetchFile(url string, fileName string) error {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer resp.Body.Close()

	if resp.Status == "200 OK" {

		file, err := os.Create(storagePath + "/" + fileName)
		if err != nil {
			fmt.Println(err)
			return err
		}
		defer file.Close()

		_, err = io.Copy(file, resp.Body)
		if err != nil {
			fmt.Println(err)
			return err
		}

	} else {
		fmt.Println(resp.Status)
	}
	return nil
}
