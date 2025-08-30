// 方向性を持ったチャネルの例

package main

import (
	"fmt"
	"time"
)

func receiver(messages <-chan int) { // 受信専用チャネルを宣言
	for {
		msg := <-messages // チャネルからメッセージを受信
		fmt.Println(time.Now().Format("15:04:05"), "Received:", msg)
	}
}

func sender(messages chan<- int) { // 送信専用チャネルを宣言
	for i := 1; ; i++ {
		fmt.Println(time.Now().Format("15:04:05"), "Sending:", i)
		messages <- i // 1秒ごとにチャネルにメッセージを送信
		time.Sleep(1 * time.Second)
	}
}

func main() {
	msgChannel := make(chan int)
	go receiver(msgChannel)
	go sender(msgChannel)
	time.Sleep(5 * time.Second)
}
