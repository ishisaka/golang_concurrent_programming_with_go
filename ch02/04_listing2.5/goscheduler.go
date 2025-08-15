// goscheduler.go
// Goスケジューラーを呼び出している例
package main

import (
	"fmt"
	"runtime"
)

/*
Note: this program is an example of what not to do; using go scheduler
to synchronize executions
*/
func sayHello() {
	fmt.Println("Hello")
}

func main() {
	go sayHello()
	// Goスケジューラーを呼び出すことでmain意外のゴルーチンに実行の機会を与える
	runtime.Gosched()
	fmt.Println("Finished")
}
