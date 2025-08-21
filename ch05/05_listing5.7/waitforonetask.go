// waitingforone.go
// 条件変数を使用してmain関数で他のゴルーチンの終了を待つ例
package main

import (
	"fmt"
	"sync"
)

func doWork(cond *sync.Cond) {
	fmt.Println("Work started")
	fmt.Println("Work finished")
	// 以下のようにすることで、main()でゴルーチンが待機状態になることが補償される。
	cond.L.Lock()   // シグナルを送信する前にミューテックスをロック
	cond.Signal()   // 条件変数に対してシグナルを送信
	cond.L.Unlock() // シグナル送信後にミューテックスをアンロック
}

func main() {
	cond := sync.NewCond(&sync.Mutex{})
	cond.L.Lock()
	for i := 0; i < 50000; i++ {
		go doWork(cond)
		fmt.Println("Waiting for child goroutine")
		cond.Wait() // doWorkの終了を待機
		fmt.Println("Child goroutine finished")
	}
	cond.L.Unlock()
}
