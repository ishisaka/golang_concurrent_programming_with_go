package main

import (
	"fmt"
	"net"

	"github.com/cutajarj/ConcurrentProgrammingWithGo/chapter10/listing10.9"
)

func main() {
	incomingConnections := make(chan net.Conn, 10)
	listing10_9.StartHttpWorkers(3, incomingConnections)
	server, _ := net.Listen("tcp", "localhost:8080")
	defer server.Close()
	for {
		conn, _ := server.Accept()
		select {
		case incomingConnections <- conn:
		default:
			fmt.Println("Server is busy")
			conn.Write([]byte("HTTP/1.1 503 Service Unavailable\r\n\r\n" +
				"<html>Busy</html>\n"))
			conn.Close()
		}
	}
}
