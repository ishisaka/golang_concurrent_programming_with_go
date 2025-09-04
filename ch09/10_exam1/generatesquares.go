// 練習問題1

package main

import "fmt"

func GenerateSquares(quit <-chan int) <-chan int {
	squares := make(chan int)
	go func() {
		i := 0
		defer close(squares)
		for {
			i += 1
			select {
			case squares <- i * i:
			case <-quit:
				return
			}
		}
	}()
	return squares
}

func main() {
	quit := make(chan int) // quitチャネルを作成
	defer close(quit)
	result := GenerateSquares(quit) // 出力チャネルでURLを返すゴルーチン関数の呼び出し
	for num := range result {       // 出力チャネルから全てのメッセージを読み込む
		fmt.Println(num) // 結果を出力
		if num >= 100 {  // 100になったら終了
			quit <- 0
		}
	}
}
