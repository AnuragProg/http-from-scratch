package main

import (
	"fmt"
	"io"

	// Uncomment this block to pass the first stage
	"net"
	"os"
	"strings"
)


func handleConn(conn net.Conn){

	data := make([]byte, 1024)
	conn.Read(data)

	route := strings.Split(strings.Split(string(data), "\r\n")[0], " ")[1]

	var response []byte
	switch {

	case strings.Split(strings.Split(string(data), "\r\n")[0], " ")[0] == "POST":
		parsedData := strings.Split(string(data), "\r\n\r\n")
		body := strings.Trim(parsedData[1], "\x00")

		filename := strings.TrimPrefix(route, "/files/")
		directory := os.Args[2]
		file, _ := os.Create(directory+"/"+filename)
		defer file.Close()
		reader := strings.NewReader(body)
		io.Copy(file, reader)
		response = []byte("HTTP/1.1 201 Created\r\n\r\n")

	case route == "/":
		response = []byte("HTTP/1.1 200 OK\r\n\r\n")

	case strings.HasPrefix(route, "/echo"):
		echoData := strings.TrimPrefix(route, "/echo/")
		response = []byte(fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Encoding: gzip\r\nContent-Type: text/plain\r\nContent-Length: %v\r\n\r\n%v", len(echoData), echoData))
	case strings.HasPrefix(route, "/user-agent"):
		headers := strings.Split(string(data), "\r\n")[1:]
		for _, header := range headers {
			keyAndValue := strings.Split(header, ": ")
			if keyAndValue[0] == "User-Agent" {
				response = []byte(fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %v\r\n\r\n%v", len(keyAndValue[1]), keyAndValue[1]))
			}
		}
	case strings.HasPrefix(route, "/files"):
		filename := strings.TrimPrefix(route, "/files/")
		directory := os.Args[2]
		file, err := os.Open(directory+"/"+filename)
		if err != nil {
			response = []byte("HTTP/1.1 404 Not Found\r\n\r\n")
			break
		}
		defer file.Close()
		
		fileStats, _ := file.Stat()
		fileData := make([]byte, fileStats.Size())
		file.Read(fileData)

		response = []byte(fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: application/octet-stream\r\nContent-Length: %v\r\n\r\n%v", len(fileData), string(fileData)))
	default:
		response = []byte("HTTP/1.1 404 Not Found\r\n\r\n")
	}
	conn.Write(response)
}

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	// Uncomment this block to pass the first stage

	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		go handleConn(conn)
	}
}
