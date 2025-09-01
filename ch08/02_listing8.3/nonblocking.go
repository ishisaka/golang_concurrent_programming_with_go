// チャネルからのノンブロッキング読み込み

package main

import (
	"fmt"
	"time"
)

func sendMsgAfter(seconds time.Duration) <-chan string {
	messages := make(chan string)
	go func() {
		time.Sleep(seconds)
		messages <- "Hello"
	}()
	return messages
}

func main() {
	// 3秒後にメッセージが送信されるチャネル
	messages := sendMsgAfter(3 * time.Second)
	for {
		select {
		// メッセージがあればチャネルから読み込む
		case msg := <-messages:
			fmt.Println("Message Received", msg)
			return // メッセージを受信したら実行を終了
		default: // メッセージが無ければdefaultケースを実行
			fmt.Println("No message waiting")
			time.Sleep(1 * time.Second)
		}
	}
}
