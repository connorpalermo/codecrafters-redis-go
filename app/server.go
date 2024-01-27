package main

import (
	"fmt"
	"io"
	"strings"

	"net"
	"os"
)

const bufferSize = 1024

var db = make(map[string]string)

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
	buffer := make([]byte, bufferSize)
	defer conn.Close()
	for {
		command, err := conn.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println("Error reading connection: ", err.Error())
			os.Exit(1)
		}
		processCommand(string(buffer[:command]), conn)
	}
}

func processCommand(message string, conn net.Conn) {
	commands := strings.Split(message, "\r\n")
	command := commands[2]
	response := "Not yet implemented"
	switch {
	case strings.EqualFold(command, "PING"):
		response = "+PONG\r\n"
	case strings.EqualFold(command, "ECHO"):
		response = "+" + commands[4] + "\r\n"
	case strings.EqualFold(command, "SET"):
		setDBValue(commands[4], commands[6])
	case strings.EqualFold(command, "GET"):
		response = retrieveDBValue(commands[4])
	default:
		fmt.Println("Command not yet implemented, ignoring for now.")
	}
	_, err := conn.Write([]byte(response))
	if err != nil {
		fmt.Println("Error writing to connection: ", err.Error())
	}
}

func setDBValue(key string, value string) {
	db[key] = value
}

func retrieveDBValue(key string) string {
	val, ok := db[key]

	if ok {
		return val
	}

	return ""
}
