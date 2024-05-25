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

	route := strings.Split(strings.Split(string(data), "\r\n")[0], " ")[1]

	var response []byte
	switch {
		case route == "/":
			response = []byte("HTTP/1.1 200 OK\r\n\r\n")
		case strings.HasPrefix(route, "/echo"):
			echoData := strings.TrimPrefix(route, "/echo/")
			response = []byte(fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %v\r\n\r\n%v", len(echoData), echoData))
		default:
			response = []byte("HTTP/1.1 404 Not Found\r\n\r\n")
	}
	conn.Write(response)
}
