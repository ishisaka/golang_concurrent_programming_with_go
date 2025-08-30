// チャネルがバッファを持つ例
// 1秒ごとにチャネルからメッセージを読み込む

package main

import (
	"fmt"
	"sync"
	"time"
)

func receiver(messages chan int, wGroup *sync.WaitGroup) {
	msg := 0
	for msg != -1 { // -1を読み込むまでループを続ける
		time.Sleep(1 * time.Second) // 1秒待つ
		msg = <-messages            // チャネルからメッセージを読み込む
		fmt.Println(time.Now().Format("15:04:05"), "Received:", msg)
	}
	wGroup.Done()
}

func main() {
	msgChannel := make(chan int, 3)  // 3つの容量をバッファとして持つチャネルを作成
	wGroup := sync.WaitGroup{}       // ウェイトグループを作成
	wGroup.Add(1)                    // ウェイトグループのサイズを1に設定する
	go receiver(msgChannel, &wGroup) // バッファ有りチャネルとウェイトグループを引数に渡してゴルーチンを起動
	for i := 1; i <= 6; i++ {
		size := len(msgChannel)
		fmt.Printf("%s Sending: %d. Buffer Size: %d\n",
			time.Now().Format("15:04:05"), i, size)
		msgChannel <- i // 1から6までの整数をメッセージとして送信
	}
	msgChannel <- -1 // -1のメッセージを送信
	wGroup.Wait()    // 受信側のゴルーチンの終了をウェイトグループで待つ
}
