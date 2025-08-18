// countdown.go
// 2つのゴルーチン間でのメモリ共有の例
package main

import (
	"fmt"
	"time"
)

func countdown(seconds *int) {
	for *seconds > 0 {
		time.Sleep(1 * time.Second)
		*seconds -= 1
	}
}

func main() {
	count := 5
	// ポインタを使用してゴルーチン間で変数を共有
	go countdown(&count)
	for count > 0 {
		time.Sleep(500 * time.Millisecond)
		fmt.Println(count)
	}
}
