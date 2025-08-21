// stingyspendycond.go
// sybc.Condによる条件変数の使用
package main

import (
	"fmt"
	"os"
	"sync"
	"time"
)

func main() {
	money := 100
	mutex := sync.Mutex{}        // ミューテックスの作成
	cond := sync.NewCond(&mutex) // ミューテックスを使い条件変数を作成
	go stingy(&money, cond)      // 両方のゴルーチンに条件変数を渡す
	go spendy(&money, cond)
	time.Sleep(2 * time.Second)
	mutex.Lock()
	fmt.Println("Money in bank account: ", money)
	mutex.Unlock()
}

func stingy(money *int, cond *sync.Cond) {
	for i := 0; i < 1000000; i++ {
		cond.L.Lock() // 条件変数上のミューテックスをロック
		*money += 10
		cond.Signal()   // money共有変数へ追加することを条件変数へシグナルを送信
		cond.L.Unlock() // 条件変数上のミューテックスをアンロック
	}
	fmt.Println("Stingy Done")
}

func spendy(money *int, cond *sync.Cond) {
	for i := 0; i < 200000; i++ {
		cond.L.Lock() // 条件変数上のミューテックスをロック
		// 残金が少ない場合は待機する
		for *money < 50 {
			cond.Wait() // ミューテックスを開放して実行を一時停止
		}
		// Wait()から戻りミューテックスを再構築し、50ドル使用する
		*money -= 50
		if *money < 0 {
			fmt.Println("Money is negative!")
			os.Exit(1)
		}
		cond.L.Unlock()
	}
	fmt.Println("Spendy Done")
}
