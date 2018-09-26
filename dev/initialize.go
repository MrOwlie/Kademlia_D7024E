package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"./kademlia"
	"./network"
)

var storagePath string = "What ever the storage path is" //TODO fix this

func main() {

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

		if !validIP4(ip) {
			fmt.Printf("%q is not a valid ip-address, exiting.", ip)
			return
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

	os.Mkdir(storagePath, 0766)

	var kadem = kademlia.GetInstance()

	network.SetPort(iOwnPort)
	network.SetHandler(kadem)
	net := network.GetInstance()
	kadem.SetNetworkHandler(net)

	var wgl sync.WaitGroup
	wgl.Add(1)
	go net.Listen(&wgl)
	wgl.Wait()

	if performJoin && !kadem.Join(ip, iPort) {
		return
	}
	if performJoin {
		kadem.IdleBucketReExploration()
	}

	var action, param1 string
	for {

		fmt.Println("\n\nPlease enter a command, type '?' for help:")
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
			_, err := ioutil.ReadFile(param1)
			if err != nil {
				fmt.Println("An error occured while reading the file!")
			} else {
				kadem.StoreFile("hej")
			}
		case action == "fetch":
			kadem.LookupData(param1)
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
