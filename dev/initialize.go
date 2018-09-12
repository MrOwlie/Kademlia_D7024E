package main

import (
	//"./d7024e"
	"os"
	"fmt"
	"strings"
	"regexp"
	"strconv"
	"./network"
	"./kademlia"
)


func main(){

	args := os.Args[1:]
	ip := args[0]
	port := args[1]
	ownPort := args[2]

	if !validIP4(ip){
		fmt.Printf("%q is not a valid ip-address, exiting.", ip)
		return
	}
	if !validPort(port){
		fmt.Printf("%q is not a valid port number, exiting.", port)
		return
	}
	if !validPort(ownPort){
		fmt.Printf("%q is not a valid port number, exiting.", port)
		return
	}

	go network.Listen("vad ska va har?", ownPort)
	go kademlia.GetInstance().Join(ip, port)

	var inp string
	for{

		fmt.Println("Please enter a command, type '?' for help:")
		fmt.Scan(&inp)

		switch inp {
		case "?":
			fmt.Println("Available commands:")
			fmt.Println("\"store 'filepath'\" to store a file.")
			fmt.Println("\"fetch 'key-value'\" to fetch a file.")
			fmt.Println("\"exit\" to exit the application.")
		case "exit":
			return
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

func validPort(port string) bool{
	if pa, err := strconv.ParseInt(port,10,64); err == nil && pa > 0 {
		return true
	}
	return false
}