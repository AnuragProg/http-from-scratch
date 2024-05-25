package main

import (
	"fmt"
	// Uncomment this block to pass the first stage
	"net"
	"strings"
	"os"
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	// Uncomment this block to pass the first stage

	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	conn, err := l.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}

	
	data := make([]byte, 1024)
	conn.Read(data)

	headers :=strings.Split(string(data), "\r\n") 
	requestInfo := strings.Split(headers[0], " ")
	fmt.Println(requestInfo)
	fmt.Println(string(data))

	var response []byte
	switch requestInfo[1]{
		case "/":
			response = []byte("HTTP/1.1 200 OK\r\n\r\n")
		default:
			response = []byte("HTTP/1.1 404 Not Found\r\n\r\n")
	}
	conn.Write(response)
}
