// cpucheck.go
// Goのランタイムが認識しているCPUの数とカーネルレベルスレッドの数を表示する
package main

import (
	"fmt"
	"runtime"
)

func main() {
	fmt.Println("Number of CPUs:", runtime.NumCPU())
	// GOMAXPROCSの値がカーネルレベルスレッド都市用意されているスレッド数
	// 通常は動作している論理CPUのか図と同じ
	fmt.Println("GOMAXPROCS:", runtime.GOMAXPROCS(0))
}
