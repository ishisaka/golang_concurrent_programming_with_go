// チャネルから関数の戻り値を取得する
// C#のasync/awaitのようなことができる

package main

import "fmt"

// findFactors は指定された整数の正の因数をすべて計算してスライスで返します。
func findFactors(number int) []int {
	result := make([]int, 0)
	for i := 1; i <= number; i++ {
		if number%i == 0 {
			result = append(result, i)
		}
	}
	return result
}

func main() {
	resultCh := make(chan []int) // intのスライスの新たなチャネルを作成
	go func() {
		// 無名関数のゴルーチン内で関数を呼び出して、結果をチャネルに書き込む
		resultCh <- findFactors(3419110721)
	}()
	fmt.Println(findFactors(4033836233)) // 普通に同期的に関数を呼び出す
	fmt.Println(<-resultCh)              // チャネルから結果を収集（読み込めるまで処理が止まる）
}
