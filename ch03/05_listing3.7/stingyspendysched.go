// stingyspendysched.go
// runtime.Gosched()でゴルーチンを強制的に実行させようとした例
// しかし、これでも上手く行かない
package main

import (
	"fmt"
	"runtime"
	"time"
)

func stingy(money *int) {
	for i := 0; i < 1000000; i++ {
		*money += 10
		runtime.Gosched() // 強制的にゴルーチンをスケジューリング
	}
	fmt.Println("Stingy Done")
}

func spendy(money *int) {
	for i := 0; i < 1000000; i++ {
		*money -= 10
		runtime.Gosched() // 強制的にゴルーチンをスケジューリング
	}
	fmt.Println("Spendy Done")
}

func main() {
	money := 100
	go stingy(&money)
	go spendy(&money)
	time.Sleep(2 * time.Second)
	fmt.Print("Money in bank account: ", money)
}

// 競合状態は減るものの、発生を抑制できているわけではない。
