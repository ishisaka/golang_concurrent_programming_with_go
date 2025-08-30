// チャネルの受信側ゴルーチンが存在しない場合の例
// 実行時ランタイムエラーとなります

package main

import (
	"fmt"
	"time"
)

func main() {
	msgChannel := make(chan string)
	go receiver(msgChannel)
	fmt.Println("Sending HELLO...")
	msgChannel <- "HELLO"
	fmt.Println("Sending THERE...")
	msgChannel <- "THERE"
	fmt.Println("Sending STOP...")
	msgChannel <- "STOP"
}

func receiver(messages chan string) {
	// 受信していない
	time.Sleep(5 * time.Second)
	fmt.Println("Receiver slept for 5 seconds")
}
