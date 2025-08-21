// waitforgroup.go
// ウェイトグループの使用例
package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func doWork(id int, wg *sync.WaitGroup) {
	i := rand.Intn(5)
	// ランダムな時間スリープ（最長5秒）
	time.Sleep(time.Duration(i) * time.Second)
	fmt.Println(id, "Done working after", i, "seconds")
	// タスクの完了を通知
	wg.Done()
}

func main() {
	wg := sync.WaitGroup{} // ウェイトグループの作成
	wg.Add(4)              // 4つのタスクがあるので、ウェイトグループのカウンタに4を加算
	for i := 1; i <= 4; i++ {
		// ウェイトグループへのポインタを渡してゴルーチンを起動
		go doWork(i, &wg)
	}
	// タスクが完了するのを待機
	wg.Wait()
	fmt.Println("All complete")
}
