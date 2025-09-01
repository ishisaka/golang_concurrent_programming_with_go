// selectを使って複数チャネルから読み込む

package main

import "time"

// writeEvery は、指定された間隔でメッセージを送信するチャネルを返します。
// msg は送信する文字列、seconds は送信の間隔を指定します。
// 戻り値はメッセージを送信するチャネルです。
func writeEvery(msg string, seconds time.Duration) <-chan string {
	messages := make(chan string) // チャネルを作成
	go func() {                   // 無名関数のゴルーチンを作成
		for {
			time.Sleep(seconds) // 指定された時間スリープ
			messages <- msg     // メッセージをチャネルへ送信
		}
	}()
	return messages // メッセージのチャネルを返す
}

func main() {
	messagesFromA := writeEvery("Tick", 1*time.Second) // チャネルAにメッセージを送信するゴルーチンを作成
	messagesFromB := writeEvery("Tock", 3*time.Second) // チャネルBにメッセージを送信するゴルーチンを作成
	for {                                              // 無限ループ
		select { // selectでチャネルをグループ化
		case msg := <-messagesFromA: // 利用可能であればチャネルAからのメッセージを出力
			println(msg)
		case msg := <-messagesFromB: // 利用可能であればチャネルBからのメッセージを出力
			println(msg)
		}
	}
}
