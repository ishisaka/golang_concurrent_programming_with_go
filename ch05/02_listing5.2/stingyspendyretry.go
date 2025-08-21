// スリープで条件変数的な機能を実装するあまりよろしくない例
package main

import (
	"fmt"
	"os"
	"sync"
	"time"
)

func stingy(money *int, mutex *sync.Mutex) {
	for i := 0; i < 1000000; i++ {
		mutex.Lock()
		*money += 10
		mutex.Unlock()
	}
	fmt.Println("Stingy Done")
}

func spendy(money *int, mutex *sync.Mutex) {
	for i := 0; i < 200000; i++ {
		mutex.Lock()
		// 残金が50ドル以下になったら稼ぐ機会を増やす
		for *money < 50 {
			// ミューテックスを開放し、他のゴルーチンによるmoneyへのアクセスを許可
			mutex.Unlock()
			time.Sleep(10 * time.Millisecond) // しばらくの間スリープ
			mutex.Lock()                      // ミューテックスを再びロック
		}
		// 十分なお金があれば再試行を続ける
		*money -= 50
		if *money < 0 {
			fmt.Println("Money is negative!")
			os.Exit(1)
		}
		mutex.Unlock()
	}
	fmt.Println("Spendy Done")
}

func main() {
	money := 100
	mutex := sync.Mutex{}
	go stingy(&money, &mutex)
	go spendy(&money, &mutex)
	time.Sleep(2 * time.Second)
	mutex.Lock()
	fmt.Println("Money in bank account: ", money)
	mutex.Unlock()
}
