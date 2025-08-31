// 練習問題3

package main

import (
	"fmt"
	"math/rand"
)

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
	resultCh := make([]chan []int, 10) // intのスライスの新たなチャネルのスライスを作成
	for i := range 10 {
		r := rand.Int31()
		resultCh[i] = make(chan []int) // チャネルのスライスにチャネルを追加
		go func(n int) {
			resultCh[n] <- findFactors(int(r))
		}(i)
	}
	for i := range 10 {
		fmt.Println(<-resultCh[i])
	}
}
