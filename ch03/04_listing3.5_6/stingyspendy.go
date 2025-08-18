// stingyspendy.go
// 競合状態が発生している例
package main

import (
	"fmt"
	"time"
)

// この関数は銀行口座の合計を保持する変数へのポインタを受け取る
func stingy(money *int) {
	for i := 0; i < 1000000; i++ {
		*money += 10 // 10ドルを家産
	}
	fmt.Println("Stingy Done")
}

// この関数は銀行口座の合計を保持する変数へのポインタを受け取る
func spendy(money *int) {
	for i := 0; i < 1000000; i++ {
		*money -= 10 // 10ドルを減算
	}
	fmt.Println("Spendy Done")
}

func main() {
	money := 100                // 初期値は100ドル
	go stingy(&money)           // ゴルーチンを起動しmoney変数へのポインタを渡す
	go spendy(&money)           // 同上
	time.Sleep(2 * time.Second) // ゴルーチンの終了を待つため2秒待つ
	fmt.Println("Money in bank account: ", money)
}

// stingyとspendyで競合状態が発生しているため結果は不確定である
