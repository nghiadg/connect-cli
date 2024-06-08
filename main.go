package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

func receiveMessage(server *net.UDPConn) {
	buffer := make([]byte, 1024)
	n, senderAddr, err := server.ReadFromUDP(buffer)
	if err != nil {
		fmt.Printf("There was an error while trying to read message from %s\n", senderAddr.String())
		os.Exit(-1)
	}

	msg := string(buffer[:n])

	fmt.Printf("-> %s: %s", senderAddr.String(), msg)

}

func sendMessage(server *net.UDPConn, reader *bufio.Reader) {
	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Opp! Can not read the message you entered")
			continue
		}

		// command -> st <address> <port> <message>. Example st 127.0.0.1 8080 hello world
		commandParts := strings.SplitN(msg, " ", 4)
		if len(commandParts) != 4 || commandParts[0] != "st" {
			fmt.Println("Invalid command format. Use: st <address> <port> <message>")
			continue
		}

		ip := commandParts[1]
		port, err := strconv.Atoi(commandParts[2])

		if err != nil {
			fmt.Println("Maybe invalid port format.")
			continue
		}

		message := commandParts[3]

		receiverAddr := net.UDPAddr{
			IP:   net.ParseIP(ip),
			Port: port,
		}

		_, err = server.WriteToUDP([]byte(message), &receiverAddr)
		if err != nil {
			fmt.Printf("Send message to %s error. Try again", receiverAddr.String())
			continue
		}

		fmt.Printf("Sen to %s: %s", receiverAddr.String(), message)
	}

}

func main() {
	reader := bufio.NewReader(os.Stdin)
	// Stdin PORT
	fmt.Print("Enter your port listener: ")
	portStr, err := reader.ReadString('\n')
	// Remove \n
	portStr = portStr[:len(portStr)-1]

	if err != nil {
		fmt.Println("Opp! Can not read the port you entered", err)
		os.Exit(-1)
	}

	port, err := strconv.Atoi(portStr)

	if err != nil {
		fmt.Println("Opp! Can not read the port you entered", err)
		os.Exit(-1)
	}

	// Register socket listener
	addr := net.UDPAddr{
		IP:   net.ParseIP("127.0.0.1"), // Listen all address
		Port: port,
	}

	server, err := net.ListenUDP("udp", &addr)
	if err != nil {
		fmt.Printf("Opp! Can not start listen on port %v\n", port)
		os.Exit(-1)
	}

	fmt.Printf("Start server success! We are listening at %s\n", addr.String())

	// goroutine for handle send message func
	go sendMessage(server, reader)

	for {
		receiveMessage(server)
	}
}
