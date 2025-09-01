// タイムアウトを伴うチャネル待ち実装

package main

import (
	"fmt"
	"os"
	"strconv"
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
	// プログラムの引数からタイムアウト値を読み込む
	t, _ := strconv.Atoi(os.Args[1])
	// 返したチャネルに3秒後にメッセージを送信するゴルーチンを開始
	messages := sendMsgAfter(3 * time.Second)
	timeoutDuration := time.Duration(t) * time.Second
	fmt.Printf("waiting for message for %d seconds...\n", t)
	select {
	case msg := <-messages:
		// チャネルにメッセージがあればそれを読み込む
		fmt.Println("Message Received", msg)
	case tNow := <-time.After(timeoutDuration):
		// チャネルとタイマーを作成し、指定された時間後にメッセージを受信
		fmt.Println("Timeout. Waited until:", tNow.Format("15:04:05"))
	}
}
