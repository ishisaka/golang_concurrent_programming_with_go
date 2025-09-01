// selectでチャネルに書き込む

package main

import (
	"fmt"
	"math"
	"math/rand"
)

// primeOnlyは整数の入力チャネルから素数のみをフィルタリングして出力チャネルに送る関数です。
func primeOnly(input <-chan int) <-chan int {
	result := make(chan int)
	go func() { // 素数のみをフィルタリングする無名のゴルーチンを生成
		for c := range input {
			isPrime := c != 1 // 1は素数ではないためcが1ではないことを確認
			// cが2からcの平方根の間までの範囲に因数を持つか確認
			for i := 2; i <= int(math.Sqrt(float64(c))); i++ {
				if c%i == 0 {
					isPrime = false
					break
				}
			}
			if isPrime {
				result <- c // cが素数であればresultチャネルへcを送信
			}
		}
	}()
	return result
}

func main() {
	numbersChannel := make(chan int)
	primes := primeOnly(numbersChannel)
	for i := 0; i < 100; { // 100個の素数を収集するまで繰り返す
		select {
		case numbersChannel <- rand.Intn(10000000000) + 1: // 1から10億までの乱数をnumbersChannelに書き込む
		case p := <-primes:
			fmt.Println("Found Prime:", p)
			i++
		}
	}
}
