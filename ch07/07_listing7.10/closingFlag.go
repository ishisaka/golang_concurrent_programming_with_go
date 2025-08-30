// チャネルのクローズとクローズの通知

package main

import (
	"fmt"
	"time"
)

func receiver(messages <-chan int) {
	for {
		// メッセージとオープンチャネルフラグを読み込む
		msg, more := <-messages
		fmt.Println(time.Now().Format("15:04:05"), "Received:", msg, more)
		time.Sleep(1 * time.Second)
		// フラグがfalseなら処理を停止する
		if !more {
			return
		}
	}
}

func main() {
	msgChannel := make(chan int)
	go receiver(msgChannel)
	for i := 1; i <= 3; i++ {
		fmt.Println(time.Now().Format("15:04:05"), "Sending:", i)
		msgChannel <- i
		time.Sleep(1 * time.Second)
	}

	close(msgChannel)
	time.Sleep(3 * time.Second)
}
