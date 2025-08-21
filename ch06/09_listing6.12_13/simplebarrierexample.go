// バリアの使用例

package main

import (
	"fmt"
	"time"

	"github.com/cutajarj/ConcurrentProgrammingWithGo/chapter6/listing6.10"
)

func workAndWait(name string, timeToWork int, barrier *listing6_10.Barrier) {
	start := time.Now()
	for {
		fmt.Println(time.Since(start), name, "is running")
		// 指定された秒数だけ作業している状況をシミュレーション
		time.Sleep(time.Duration(timeToWork) * time.Second)
		fmt.Println(time.Since(start), name, "is waiting on barrier")
		barrier.Wait() // 他のゴルーチンが追いつくのを待つ
	}
}

func main() {
	// 2つのゴルーチンが参加する新しいバリアを作成
	barrier := listing6_10.NewBarrier(2)
	// ゴルーチンを開始
	go workAndWait("Red", 4, barrier)
	go workAndWait("Blue", 10, barrier)
	// 100秒間スリープ
	time.Sleep(100 * time.Second)
}
