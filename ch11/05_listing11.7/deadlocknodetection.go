// パニックは起きなくなったが、まだデッドロックの検出は実装されていないコード

package main

import (
	"fmt"
	"sync"
	"time"
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
	// メッセージを表示する前にウェイトグループで待つゴルーチンを作成
	go func() {
		wg.Wait()
		fmt.Println("Done waiting on waitgroup")
		// このゴルーチンが永遠にwgの終了を待ってしまうので、ランタイムがデッドロックを検出できない。
	}()
	time.Sleep(30 * time.Second) // 30秒待つ
	fmt.Println("Done")          // メッセージを表示して、プログラムは強制終了
}
