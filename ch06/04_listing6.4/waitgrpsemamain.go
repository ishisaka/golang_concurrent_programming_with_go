// 自作のウェイトグループの使用例

package main

import (
	"fmt"

	"github.com/cutajarj/ConcurrentProgrammingWithGo/chapter6/listing6.3"
)

func doWork(id int, wg *listing6_3.WaitGrp) {
	fmt.Println(id, "Done working ")
	wg.Done() // ゴルーチンの終了をDone()を呼び出して通知する
}

func main() {
	wg := listing6_3.NewWaitGrp(4) // サイズ4のウェイトグループを作成
	for i := 1; i <= 4; i++ {
		go doWork(i, wg) // ウェイトグループへのポインタを渡してゴルーチンを作成
	}
	wg.Wait() // 作業の完了をウェイトグループを使用して待つ
	fmt.Println("All complete")
}
