// signalbeforewait.go
// 条件変数を使用してmain関数で他のゴルーチンの終了を待ちたいが、
// 期待通りに動作せず、エラーとなる例
package main

import (
	"fmt"
	"sync"
)

/*
注：このプログラムにはデモ目的のバグが含まれています
次のリストで、この問題を解決する方法を説明します
*/
func doWork(cond *sync.Cond) {
	fmt.Println("Work started")
	fmt.Println("Work finished")
	cond.Signal()
}

func main() {
	cond := sync.NewCond(&sync.Mutex{})
	cond.L.Lock()
	for i := 0; i < 50000; i++ { // 5000回繰り返す
		go doWork(cond) // ゴルーチンを起動し何らかの作業をシミュレートする
		fmt.Println("Waiting for child goroutine")
		cond.Wait() // ゴルーチンからの終了シグナルを待つ
		fmt.Println("Child goroutine finished")
	}
	cond.L.Unlock()
}
