// nonblockingmutex.go
// ノンブロッキング・ミューテックス・ロックの例
package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/cutajarj/ConcurrentProgrammingWithGo/chapter4/listing4.5"
)

func main() {
	mutex := sync.Mutex{}
	var frequency = make([]int, 26)
	for i := 2000; i <= 2200; i++ {
		url := fmt.Sprintf("https://rfc-editor.org/rfc/rfc%d.txt", i)
		go listing4_5.CountLetters(url, frequency, &mutex)
	}
	for i := 0; i < 100; i++ {
		time.Sleep(100 * time.Millisecond) // 100ミリ秒のスリープ
		if mutex.TryLock() {               // ミューテックスの獲得を試みる
			// ミューテックスが獲得できたら頻度カウンターを表示して、ミューテックスを開放する。
			for i, c := range listing4_5.AllLetters {
				fmt.Printf("%c-%d ", c, frequency[i])
			}
			fmt.Println()
			mutex.Unlock()
		} else {
			// ミューテックスが獲得できない場合にはメッセージを表示して、再試行する。
			fmt.Println("Mutex already being used")
		}
	}
}
