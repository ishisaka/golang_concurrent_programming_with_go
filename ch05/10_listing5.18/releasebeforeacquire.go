// セマフォの使用例
package main

import (
	"fmt"

	"github.com/cutajarj/ConcurrentProgrammingWithGo/chapter5/listing5.16"
)

func main() {
	semaphore := listing5_16.NewSemaphore(0)
	for i := 0; i < 50000; i++ { // 5000回繰り返す
		go doWork(semaphore) // セマフォへのポインタを渡してゴルーチンを開始
		fmt.Println("Waiting for child goroutine")
		// タスク完了を示すセマフォ内の許可が利用可能になるまで待機
		semaphore.Acquire()
		fmt.Println("Child goroutine finished")
	}
}

func doWork(semaphore *listing5_16.Semaphore) {
	fmt.Println("Work started")
	fmt.Println("Work finished")
	// ゴルーチンの終了時に、main()ゴルーチンへ通知するために許可を開放
	semaphore.Release()
}
