// for rangeを使ったチャネルからのメッセージ読み込みの例

package main

import (
	"fmt"
	"time"
)

func receiver(messages <-chan int) {
	// チャネルがクローズされるまでメッセージを読み込む
	for msg := range messages {
		fmt.Println(time.Now().Format("15:32:06"), "Received:", msg)
		time.Sleep(1 * time.Second)
	}
	fmt.Println("Receiver done")
}

func main() {
	msgChannel := make(chan int)
	go receiver(msgChannel)
	for i := 1; i <= 3; i++ {
		fmt.Println(time.Now().Format("15:32:06"), "Sending:", i)
		msgChannel <- i
		time.Sleep(1 * time.Second)
	}
	close(msgChannel)
	time.Sleep(3 * time.Second)
}
