package main

import (
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"net"
	"os"

	"github.com/hdt3213/rdb/parser"
)

const bufferSize = 1024

var db = make(map[string]string)
var ttl = make(map[string]int64)
var rdbKeys = make(map[string]string)
var properties = make(map[string]string)

func main() {

	populateProperties()

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
		db[commands[4]] = commands[6]
		if len(commands) > 8 && strings.EqualFold(commands[8], "px") {
			ttl[commands[4]] = makeTimestamp(commands[10])
		}
		response = "+OK\r\n"
	case strings.EqualFold(command, "GET") && properties["dbfilename"] == "":
		response = processGetCommand(retrieveDBValue(commands[4]))
	case strings.EqualFold(command, "CONFIG"):
		response = processConfigCommand(commands)
	case strings.EqualFold(command, "KEYS"):
		array := retrieveKeysFromFile()
		response = "*" + strconv.Itoa(len(array)) + "\r\n"
		for i := 0; i < len(array); i++ {
			response += "$" + strconv.Itoa(len(array[i][0])) + "\r\n" + array[i][0] + "\r\n"
		}
	case strings.EqualFold(command, "GET"):
		retrieveKeysFromFile()
		response = "$" + strconv.Itoa(len(rdbKeys[commands[4]])) + "\r\n" + rdbKeys[commands[4]] + "\r\n"
	default:
		fmt.Println("Command not yet implemented, ignoring for now.")
	}
	_, err := conn.Write([]byte(response))
	if err != nil {
		fmt.Println("Error writing to connection: ", err.Error())
	}
}

func makeTimestamp(milliseconds string) int64 {
	n, err := strconv.ParseInt(milliseconds, 10, 64)
	if err == nil {
		fmt.Printf("%d of type %T", n, n)
	}
	return (time.Now().UnixNano() / 1e6) + n
}

func retrieveDBValue(key string) string {
	val, ok := db[key]

	if ok {
		return val
	}

	return ""
}

func processConfigCommand(commands []string) string {
	processed := ""
	if strings.EqualFold(commands[4], "GET") {
		keyLen := len(commands[6])
		valLen := len(properties[commands[6]])
		processed = "*2\r\n$" + strconv.Itoa(keyLen) + "\r\n" + commands[6] +
			"\r\n$" + strconv.Itoa(valLen) + "\r\n" + properties[commands[6]] + "\r\n"
	}
	return processed
}

func processGetCommand(command string) string {
	response := ""
	if val, ok := ttl[command]; ok {
		if val-(time.Now().UnixNano()/1e6) > 0 {
			response = "+" + command + "\r\n"
		} else {
			response = "$-1\r\n"
		}
	} else {
		response = "+" + command + "\r\n"
	}
	return response
}

func retrieveKeysFromFile() [][]string {
	array := processRDB()
	return array
}

func populateProperties() {
	args := os.Args[1:]

	for i := 0; i < len(args); i++ {
		if i%2 == 0 {
			key := strings.ReplaceAll(args[i], "-", "")
			properties[key] = args[i+1]
		}
	}
}

func processRDB() [][]string {
	fileName := properties["dir"] + "/" + properties["dbfilename"]
	key, value := "", ""
	var keyVals [][]string
	rdbFile, err := os.Open(fileName)
	if err != nil {
		panic("open dump.rdb failed")
	}
	defer func() {
		_ = rdbFile.Close()
	}()
	decoder := parser.NewDecoder(rdbFile)
	err = decoder.Parse(func(o parser.RedisObject) bool {
		current := make([]string, 2)
		str := o.(*parser.StringObject)
		key = str.Key
		value = string(str.Value)
		current[0] = key
		current[1] = value
		rdbKeys[key] = value
		keyVals = append(keyVals, current)
		if key == "" || value == "" {
			return false
		}
		key = ""
		value = ""
		current = nil
		return true
	})
	if err != nil {
		panic(err)
	}

	return keyVals
}
