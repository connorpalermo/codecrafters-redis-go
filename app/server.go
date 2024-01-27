package main

import (
	"fmt"
	"strings"

	"net"
	"os"
)

const bufferSize = 1024

func main() {

	listener, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}
	defer listener.Close()
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}
		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	defer conn.Close()
	bytes := make([]byte, bufferSize)
	command, err := conn.Read(bytes)
	if err != nil {
		fmt.Println("Error reading: ", err.Error())
		os.Exit(1)
	}
	processCommand(string(bytes[:command]), conn)
}

func processCommand(message string, conn net.Conn) {
	commands := strings.Split(message, "\r\n")
	response := "Not yet implemented"
	switch {
	case strings.EqualFold(commands[2], "PING"):
		response = "+PONG\r\n"
	case strings.Contains(commands[2], "ECHO") || strings.Contains(commands[2], "echo"):
		response = "+" + commands[4] + "\r\n"
	default:
		fmt.Println("Command not yet implemented, ignoring for now.")
	}
	_, err := conn.Write([]byte(response))
	if err != nil {
		fmt.Println("Error writing to connection: ", err.Error())
	}
}
