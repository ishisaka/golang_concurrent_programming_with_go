// nilのチャネル
// このプログラムは実行するとデッドロックのエラーになります。

package main

import "fmt"

func main() {
	var ch chan string = nil // nilのチャネルを作成
	ch <- "Message"          // nilのチャンネルメッセージを送信しようとして実行が停止される
	fmt.Println("This line will not be printed")
}
