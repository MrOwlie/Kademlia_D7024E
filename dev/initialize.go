package main

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"./kademlia"
	"./messageBufferList"
	"./metadata"
	"./routingTable"

	"./network"
)

func main() {

	storagePath, _ := filepath.Abs("../storage")
	downloadPath, _ := filepath.Abs("../downloads")

	args := os.Args[1:]
	var ownPort, ip, port string
	var iOwnPort, iPort int
	performJoin := true
	if len(args) == 1 {
		ownPort = args[0]
		performJoin = false
	} else if len(args) == 3 {
		ip = args[0]
		port = args[1]
		ownPort = args[2]

		if !validIP4(ip) && ip != "localhost" {
			for {
				res, err := net.LookupHost(ip)
				if err != nil {
					fmt.Println(ip, " is not a valid ip-address or host-name, exiting.")
					time.Sleep(5 * time.Second)
				} else {
					ip = res[0]
					break
				}
			}
		}
		if !validPort(port) {
			fmt.Printf("%q is not a valid port number, exiting.", port)
			return
		}
		iPort, _ = strconv.Atoi(port)

	} else {
		fmt.Println("Unvalid amount of arguments")
		return
	}

	if !validPort(ownPort) {
		fmt.Printf("%q is not a valid port number, exiting.", ownPort)
		return
	}
	iOwnPort, _ = strconv.Atoi(ownPort)

	if ex, perr := pathExists(storagePath); perr != nil {
		fmt.Println("file system error: ", perr)
		return
	} else if !ex {
		os.Mkdir(storagePath, 0766)
	}

	if ex, perr := pathExists(downloadPath); perr != nil {
		fmt.Println("file system error: ", perr)
		return
	} else if !ex {
		os.Mkdir(downloadPath, 0766)
	}

	RTable := routingTable.NewRoutingTable()
	MBList := &messageBufferList.MessageBufferList{}
	MData := metadata.NewFileMetaData()
	var kadem = kademlia.NewKademliaObject(RTable, MBList, MData)

	net := network.NewNetwork(iOwnPort, "")
	net.SetHandler(kadem)
	kadem.SetNetworkHandler(net)

	var wgl sync.WaitGroup
	wgl.Add(1)
	go net.Listen(&wgl)
	go net.ListenFileServer()
	wgl.Wait()

	if performJoin && !kadem.Join(ip, iPort) {
		return
	}
	/*if performJoin {
		kadem.IdleBucketReExploration()
	}*/

	var action, param1 string
	for {

		fmt.Println("\n\nPlease enter a command, type '?' for help :D :")
		fmt.Scanln(&action, &param1)

		switch {
		case action == "?":
			fmt.Println("Available commands:")
			fmt.Println("\"store 'filepath'\" to store a file.")
			fmt.Println("\"fetch 'key-value'\" to fetch a file.")
			fmt.Println("\"exit\" to exit the application.")
		case action == "exit":
			return
		case action == "store":
			path, _ := filepath.Abs(param1)
			/*_, err := ioutil.ReadFile(path)
			if err != nil {
				fmt.Println("An error occured while reading the file!")
			} else {*/
			kadem.StoreFile(path)
			//}
		case action == "fetch":
			_, closest, _ := kadem.LookupData(param1)
			fmt.Println("Nodes found :", len(closest))
			for _, c := range closest {
				fmt.Println(c.ID)
			}
		default:
			fmt.Println("The command entered is invalid, try again.")
		}
	}

}

func validIP4(ipAddress string) bool {
	ipAddress = strings.Trim(ipAddress, " ")

	re, _ := regexp.Compile(`^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$`)
	if re.MatchString(ipAddress) {
		return true
	}
	return false
}

func validPort(port string) bool {
	if pa, err := strconv.Atoi(port); err == nil && pa > 0 {
		return true
	}
	return false
}

func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

type hubDuplex struct {
	incoming chan []string
	outgoing chan []string
}

type hub struct {
	RoutingTable     *routingTable.RoutingTable
	MBList           *messageBufferList.MessageBufferList
	MetaData         *metadata.FileMetaData
	KademliaInstance *kademlia.Kademlia
}

func newHub(port int) *hub {
	hub := &hub{}
	hub.RoutingTable = routingTable.NewRoutingTable()
	hub.MBList = &messageBufferList.MessageBufferList{}
	hub.MetaData = metadata.NewFileMetaData()
	hub.KademliaInstance = kademlia.NewKademliaObject(hub.RoutingTable, hub.MBList, hub.MetaData)

	net := network.NewNetwork(port, "")
	net.SetHandler(hub.KademliaInstance)
	hub.KademliaInstance.SetNetworkHandler(net)
}
