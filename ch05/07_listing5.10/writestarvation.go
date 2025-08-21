// writestarvation.go
// リーダーゴルーチンがリーダーロックを専有し、書き込みアクセスを待たせる
package main

import (
	"fmt"
	"time"

	"github.com/cutajarj/ConcurrentProgrammingWithGo/chapter4/listing4.12"
)

func main() {
	rwMutex := listing4_12.ReadWriteMutex{}
	for i := 0; i < 2; i++ {
		// 2つのゴルーチンを開始
		go func() {
			for { // 無限に繰り返す
				rwMutex.ReadLock()
				// リーダーロックしたまま1秒間スリープ
				time.Sleep(1 * time.Second)
				fmt.Println("Read done")
				rwMutex.ReadUnlock()
			}
		}()
	}
	time.Sleep(1 * time.Second)
	rwMutex.WriteLock()
	fmt.Println("Write finished")
}
