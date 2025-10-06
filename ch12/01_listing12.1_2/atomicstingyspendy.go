// アトミック処理の例

package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

func stingy(money *int32) {
	for range 1000000 {
		atomic.AddInt32(money, 10) // money共有変数にアトミックに10ドル加算
	}
	fmt.Println("Stingy Done")
}

func spendy(money *int32) {
	for range 1000000 {
		atomic.AddInt32(money, -10) // money共有変数からアトミックに10ドル減算
	}
	fmt.Println("Spendy Done")
}

func main() {
	money := int32(100) // 32ビット整数で100を作成
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		stingy(&money)
		wg.Done()
	}()
	go func() {
		spendy(&money)
		wg.Done()
	}()
	wg.Wait() // 両方のゴルーチンが終わるまで待機
	// 共有変数の値を読み込み、コンソールに表示
	fmt.Println("Money in account: ", atomic.LoadInt32(&money))
}
