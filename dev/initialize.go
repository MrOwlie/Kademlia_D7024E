package main

import (
	"Kademlia_D7024E/dev/kademlia"
	"Kademlia_D7024E/dev/network"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func main() {

	args := os.Args[1:]
	ip := args[0]
	port := args[1]
	ownPort := args[2]

	if !validIP4(ip) {
		fmt.Printf("%q is not a valid ip-address, exiting.", ip)
		return
	}
	if !validPort(port) {
		fmt.Printf("%q is not a valid port number, exiting.", port)
		return
	}
	if !validPort(ownPort) {
		fmt.Printf("%q is not a valid port number, exiting.", port)
		return
	}

	var kadem = kademlia.GetInstance()

	iOwnPort, _ := strconv.Atoi(ownPort)
	iPort, _ := strconv.Atoi(port)

	go network.Listen("vad ska va har?", iOwnPort)
	go kadem.Join(ip, iPort)

	var action, param1 string
	for {

		fmt.Println("Please enter a command, type '?' for help:")
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
			data, err := ioutil.ReadFile(param1)
			if err != nil {
				fmt.Println("An error occured while reading the file!")
			} else {
				kadem.Store(data)
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
