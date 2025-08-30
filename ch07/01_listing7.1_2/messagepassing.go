// チャネルの使用例

package main

import "fmt"

func main() {
	msgChannel := make(chan string) // チャネルの作成
	go receiver(msgChannel)         // 作成したチャネルを引数としてゴルーチンを起動
	fmt.Println("Sending HELLO...")
	msgChannel <- "HELLO" // チャネルメッセージを送信
	fmt.Println("Sending THERE...")
	msgChannel <- "THERE"
	fmt.Println("Sending STOP...")
	msgChannel <- "STOP"
}

func receiver(messages chan string) {
	msg := ""
	for msg != "STOP" {
		msg = <-messages // メッセージを受信
		fmt.Println("Received:", msg)
	}
}
