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
	"./network"
	"./routingTable"
	"./serverd"
)

func main() {

	storagePath, _ := filepath.Abs("../storage")
	downloadPath, _ := filepath.Abs("../downloads")

	args := os.Args[1:]
	var ownPort, ip, port string
	var iOwnPort int
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

	hub := newHub(iOwnPort)

	CLIChannels := &hubDuplex{make(chan []string), make(chan []string)}
	hub.addConnector(CLIChannels)
	fmt.Println("Listening to commands from CLI")

	APIChannels := &hubDuplex{make(chan []string), make(chan []string)}
	ApiServer := serverd.NewAPIServer(APIChannels.outgoing, APIChannels.incoming)
	go ApiServer.ListenApiServer()
	hub.addConnector(APIChannels)
	fmt.Println("Listening to commands from API")
	hub.Listen()
	if performJoin {
		CLIChannels.incoming <- []string{"join", ip, port}
		response := <-CLIChannels.outgoing
		fmt.Println(response[1])
		if response[0] != "success" {
			return
		}
	}
	/*if performJoin {
		kadem.IdleBucketReExploration()
	}*/

	var CLIExit sync.WaitGroup
	CLIExit.Add(1)

	go func() {
		var action, param1 string
		for {

			fmt.Println("\n\nPlease enter a command, type '?' for help :D :")
			fmt.Scanln(&action, &param1)

			switch {
			case action == "?":
				fmt.Println("Available commands:")
				fmt.Println("\"store 'filepath'\" to store a file.")
				fmt.Println("\"cat 'key-value'\" to fetch a file.")
				fmt.Println("\"pin\" 'key-value'\" to pin a file.")
				fmt.Println("\"unpin\" 'key-value'\" to unpin a file.")
				fmt.Println("\"exit\" to exit the application.")
			case action == "exit":
				CLIExit.Done()

			case action == "store":
				path, _ := filepath.Abs(param1)
				CLIChannels.incoming <- []string{"store", path}
				response := <-CLIChannels.outgoing
				fmt.Println(response[1])

			case action == "cat":
				CLIChannels.incoming <- []string{"cat", param1}
				response := <-CLIChannels.outgoing
				fmt.Println(response[1])

			case action == "pin":
				CLIChannels.incoming <- []string{"pin", param1}
				response := <-CLIChannels.outgoing
				fmt.Println(response[1])

			case action == "unpin":
				CLIChannels.incoming <- []string{"unpin", param1}
				response := <-CLIChannels.outgoing
				fmt.Println(response[1])

			default:
				fmt.Println("The command entered is invalid, try again.")
			}
		}
	}()

	CLIExit.Wait()
	return

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

type hubMessage struct {
	command   []string
	Connector *hubDuplex
}

type hub struct {
	RoutingTable     *routingTable.RoutingTable
	MBList           *messageBufferList.MessageBufferList
	MetaData         *metadata.FileMetaData
	KademliaInstance *kademlia.Kademlia
	Connectors       []*hubDuplex
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

	var wgl sync.WaitGroup
	wgl.Add(1)
	go net.Listen(&wgl)
	go net.ListenFileServer()
	wgl.Wait()

	return hub
}

func (h *hub) addConnector(Connector *hubDuplex) {
	h.Connectors = append(h.Connectors, Connector)
}

func (h *hub) Listen() {

	mergeChan := make(chan hubMessage)

	for i := 0; i < len(h.Connectors); i++ {
		go func(a int) {
			for {
				command := <-h.Connectors[a].incoming
				mergeChan <- hubMessage{command, h.Connectors[a]}
			}
		}(i)
	}

	go func() {
		for {
			for message := range mergeChan {

				command := message.command
				var response []string

				switch {
				case command[0] == "join":
					iPort, _ := strconv.Atoi(command[2])
					fmt.Println("Attempting to join network via ", command[1], ":", iPort)
					joined := h.KademliaInstance.Join(command[1], iPort)
					if joined {
						response = []string{"success", "Successfully joined " + command[1]}
					} else {
						response = []string{"fail", "Join attempt towards " + command[1] + " timed out"}
					}

				case command[0] == "store":
					fmt.Println("Attempting to store file on network")
					filename := h.KademliaInstance.StoreFile(command[1])
					if len(filename) > 0 {
						response = []string{"success", "Successfully stored file as " + filename}
					} else {
						response = []string{"fail", "Storing file failed"}
					}

				case command[0] == "cat":
					fmt.Println("Attempting to fetch file ", command[1])
					path, _, found := h.KademliaInstance.LookupData(command[1])
					if found {
						response = []string{"success", "Successfully downloaded file", path}
					} else {
						response = []string{"fail", "File was not found"}
					}

				case command[0] == "pin":
					fmt.Println("Attempting to pin file ", command[1])
					if h.KademliaInstance.PinFile(command[1]) {
						response = []string{"success", "Pinned " + command[1]}
					} else {
						response = []string{"fail", "File was not found"}
					}

				case command[0] == "unpin":
					fmt.Println("Attempting to unpin file ", command[1])
					if h.KademliaInstance.PinFile(command[1]) {
						response = []string{"success", "Unpinned " + command[1]}
					} else {
						response = []string{"fail", "File was not found"}
					}

				default:
					response = []string{"fail", "The command was unrecognized"}
				}

				message.Connector.outgoing <- response
			}
		}
	}()
}
