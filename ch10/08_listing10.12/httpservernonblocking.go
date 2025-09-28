package main

import (
	"fmt"
	"net"

	"github.com/cutajarj/ConcurrentProgrammingWithGo/chapter10/listing10.9"
)

func main() {
	incomingConnection := make(chan net.Conn)
	listing10_9.StartHttpWorkers(3, incomingConnection)
	server, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	defer server.Close()
	for {
		conn, err := server.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
		}
		select {
		case incomingConnection <- conn:
		// ゴルーチンが作業キューから読み込んでいない場合、defaultケースが実行される
		default:
			// クライアントに「ビジー状態」エラーメッセージを返す。
			fmt.Println("Connection channel is full. Closing.")
			conn.Write([]byte("HTTP/1.1 429 Too Many Requests\r\n\r\n" + "<html>Busy</html>"))
			conn.Close()
		}
	}

}
