package main

import (
	"net"

	"github.com/cutajarj/ConcurrentProgrammingWithGo/chapter10/listing10.9"
)

func main() {
	incomingConnections := make(chan net.Conn) // 作業キューチャネルを作成
	// 3つのゴルーチンでワーカープールを開始
	listing10_9.StartHttpWorkers(3, incomingConnections)
	// 8000ポートでTCP接続を待ち受ける
	server, _ := net.Listen("tcp", "localhost:8080")
	defer server.Close()
	for {
		conn, _ := server.Accept()  // クライアントからの新たな接続があるまで待つ
		incomingConnections <- conn // 作業キューチャネルへ接続を渡す。
	}
}
