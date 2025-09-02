// quitチャネルパターン

package main

import "fmt"

func printNumbers(numbers <-chan int, quit chan int) {
	go func() {
		for i := 0; i < 10; i++ { // numberチャネルから10個読み込む
			fmt.Println(<-numbers)
		}
		close(quit) // quitチャネルをクローズ
	}()
}

func main() {
	numbers := make(chan int)   // numbersチャネルを作成
	quit := make(chan int)      // quitチャネルを作成
	printNumbers(numbers, quit) // チャネルを渡してprintNumbersを呼び出す
	next := 0
	for i := 1; ; i++ {
		next += i // 次の三角形を作成
		select {
		case numbers <- next: // numbersチャネルへ数値を送信
		case <-quit:
			// quitチャネルの待ちが解除されたら、メッセージ出力して実行を終了。（クローズされたチャネルはゼロ値を返す。）
			fmt.Println("Quitting number generation")
			return
		}
	}
}
