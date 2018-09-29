package network

import (
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"sync"
)

const MAX_PACKET_SIZE int = 5120 //TODO Calculate actual max packet size.
var storagePath, _ = filepath.Abs("../storage")

type Handler interface {
	HandleIncomingRPC([]byte, string)
}

type network struct {
	port         int
	ip           string
	conn         *net.UDPConn
	msgHandle    Handler
	sendingMutex sync.Mutex
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

func (network *network) Listen(wg *sync.WaitGroup) {
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
	wg.Done()
	for {
		var data [MAX_PACKET_SIZE]byte
		n, addr, err := conn.ReadFromUDP(data[0:])
		fmt.Println("got msg")
		if err != nil {
			fmt.Println(err)
			continue
		}
		strAddr := addr.IP.String() + ":" + strconv.Itoa(addr.Port)
		go network.msgHandle.HandleIncomingRPC(data[0:n], strAddr)
	}

}

func (network *network) SendMessage(addr string, data *[]byte) {
	network.sendingMutex.Lock()
	defer network.sendingMutex.Unlock()

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

func (network *network) ListenFileServer() {
	http.Handle("/storage", http.FileServer(http.Dir(storagePath)))
	err := http.ListenAndServe(":"+strconv.Itoa(network.port), nil)
	if err != nil {
		fmt.Println(err)
	}
}

func (network *network) FetchFile(url string, filePath string) error {

	resp, err := http.Get("http://"+url)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer resp.Body.Close()

	if resp.Status == "200 OK" {

		file, err := os.Create(filePath)
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
		return errors.New(resp.Status)
	}
	return nil
}
