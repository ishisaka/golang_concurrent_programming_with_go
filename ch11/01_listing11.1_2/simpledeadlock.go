// デッドロックが発生する例

package main

import (
	"fmt"
	"sync"
	"time"
)

func red(lock1, lock2 *sync.Mutex) {
	for {
		fmt.Println("Red: Acquiring lock1")
		lock1.Lock()
		fmt.Println("Red: Acquiring lock2")
		lock2.Lock()
		// 両方のロックを獲得し、保持
		fmt.Println("Red: Both locks Acquired")
		lock1.Unlock()
		lock2.Unlock()
		// 両方のロックを解放
		fmt.Println("Red: Locks Released")
	}
}

func blue(lock1, lock2 *sync.Mutex) {
	for {
		fmt.Println("Blue: Acquiring lock2")
		lock2.Lock()
		fmt.Println("Blue: Acquiring lock1")
		lock1.Lock()
		// 両方のロックを獲得し、保持
		fmt.Println("Blue: Both locks Acquired")
		lock1.Unlock()
		lock2.Unlock()
		// 両方のロックを解放
		fmt.Println("Blue: Locks Released")
	}
}

func main() {
	lockA := sync.Mutex{}
	lockB := sync.Mutex{}
	go red(&lockA, &lockB)  // red()ゴルーチンを起動
	go blue(&lockA, &lockB) // blue()ゴルーチンを起動
	// 両方のゴルーチンを20秒間実行
	time.Sleep(20 * time.Second)
	fmt.Println("Done")
}
