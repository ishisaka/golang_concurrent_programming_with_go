// 素数パイプラインアルゴリズム

package main

import "fmt"

func primeMultipleFilter(numbers <-chan int, quit chan<- int) {
	var right chan int
	p := <-numbers // 入力チャネルから素数pを含む最初のメッセージを受信し表示
	fmt.Println(p)
	for n := range numbers { // 入力チャネルから次の数を読み込む
		if n%p != 0 { // 読み込んだ数がpの倍数なら破棄
			// 現在のゴルーチンが右側のチャネルを持っていなければ新たなゴルーチンを開始して、チャネルで接続
			if right == nil {
				right = make(chan int)
				go primeMultipleFilter(right, quit)
			}
			right <- n // 右側のチャネルにフィルタされた数を渡す
		}
	}
	if right == nil {
		// フィルタする数がなく、右側のチャネルの度ルーチンがいなければquitチャネルをクローズ
		close(quit)
	} else {
		// そうでなければ右側のチャネルをクローズ
		close(right)
	}
}

func main() {
	numbers := make(chan int) // 素数の倍数フィルタへ送信する入力チャネルを作成
	quit := make(chan int)    // 共通のquitチャネルを作成
	// パイプライン内の最初のゴルーチンを開始し、numberチャネルとquitチャネルを渡す
	go primeMultipleFilter(numbers, quit)
	// 入力チャネルへ2以上100000未満の連続した数を送信
	for i := 2; i < 100000; i++ {
		numbers <- i
	}
	close(numbers) // 入力チャネルをクローズし、数がこれ以上ないことを通知
	<-quit         // quitチャネルがクローズされるのを待つ
}
