package main

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"net"
	"os"
)

func main() {

	listener, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}

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
	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				fmt.Println("End of File")
				os.Exit(1)
			}
			fmt.Println(err.Error())
		}
		fmt.Println("Received message: ", line)
		command := strings.TrimSpace(strings.ToUpper(line))
		processCommand(command, writer)
	}
}

func processCommand(command string, writer *bufio.Writer) {
	if command == "PING" {
		_, err := writer.WriteString("+PONG\r\n")
		if err != nil {
			fmt.Println("Error writing to connection: ", err.Error())
		}
		writer.Flush()
	} else {
		fmt.Println("Command not yet implemented, ignoring for now.")
	}
}
