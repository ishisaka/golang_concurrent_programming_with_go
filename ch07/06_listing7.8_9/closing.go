// チャネルを閉じる例
// ただし、このコードには問題がある。

package main

import (
	"fmt"
	"time"
)

func receiver(messages <-chan int) { // 受信専用チャネルを宣言
	for {
		msg := <-messages // チャネルから1つのメッセージを読み込む
		fmt.Println(time.Now().Format("15:04:05"), "Received:", msg)
		time.Sleep(1 * time.Second) // 1秒間スリープ
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
	close(msgChannel)           // チャネルをクローズ
	time.Sleep(3 * time.Second) // このスリープの間、0が出力し続けてしまう。
}
