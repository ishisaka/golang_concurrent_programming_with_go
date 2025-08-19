// stingyspendymutex.go
// Mutexを使ったクリティカルセクション保護の例
package main

import (
	"fmt"
	"sync"
	"time"
)

func stingy(money *int, mutex *sync.Mutex) { // Mutexのポインタを受け取る
	for range 1000000 {
		mutex.Lock() // Mutexをロック
		*money += 10
		mutex.Unlock() // Mutexをアンロック
	}
	fmt.Println("Stingy Done")
}

func spendy(money *int, mutex *sync.Mutex) { // Mutexのポインタを受け取る
	for range 1000000 {
		mutex.Lock()
		*money -= 10
		mutex.Unlock()
	}
	fmt.Println("Spendy Done")
}

func main() {
	money := 100
	mutex := sync.Mutex{}     // 新たなミューテックスを作成
	go stingy(&money, &mutex) // 2つのゴルーチンへMutexのポインタを渡す
	go spendy(&money, &mutex)
	time.Sleep(2 * time.Second)
	mutex.Lock()
	fmt.Println("Money in bank account: ", money)
	mutex.Unlock()
}
