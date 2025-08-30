// 送信側が存在しない例
// これも最終的にランタイムエラーとなる

package main

import (
	"fmt"
	"time"
)

func main() {
	msgChannel := make(chan string)
	go sender(msgChannel)
	fmt.Println("Reading message from channel...")
	msg := <-msgChannel // 送信側ゴルーチンが存在しないのでここで待たされる。最終的にエラー。
	fmt.Println("Received:", msg)
}

func sender(messages chan string) {
	time.Sleep(5 * time.Second)
	fmt.Println("Sender slept for 5 seconds")
}
