// Goランタイムがデッドロックを検出してパニックになるコード例

package main

import (
	"fmt"
	"sync"
)

func lockBoth(lock1, lock2 *sync.Mutex, wg *sync.WaitGroup) {
	for i := 0; i < 1000; i++ {
		// 両方のミューテックスをロックしてアンロックする
		lock1.Lock()
		lock2.Lock()
		lock1.Unlock()
		lock2.Unlock()
	}
	wg.Done() // ウェイトグループを終了
}

func main() {
	lockA, lockB := &sync.Mutex{}, &sync.Mutex{}
	wg := sync.WaitGroup{}
	wg.Add(2)
	go lockBoth(lockA, lockB, &wg) // 同時に両方のミューテックスをロックする2つのゴルーチンを開始
	go lockBoth(lockB, lockA, &wg)
	wg.Wait() // ゴルーチンの終了を待つ
	fmt.Println("Done")
}
