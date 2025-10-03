// 少し複雑なデッドロックが発生する例

package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/cutajarj/ConcurrentProgrammingWithGo/chapter11/listing11.3_4"
)

// main は複数の銀行口座間でランダムに送金を行うゴルーチンを生成し、並行処理を実行する関数です。
// このプログラムはいずれかの段階でデッドロックが発生し、処理が完了できない。
// シングルコアの環境などではデッドロックが発生しないこともあります。
func main() {
	accounts := []listing11_3_4.BankAccount{
		*listing11_3_4.NewBankAccount("Sam"),
		*listing11_3_4.NewBankAccount("Paul"),
		*listing11_3_4.NewBankAccount("Amy"),
		*listing11_3_4.NewBankAccount("Mia"),
	}
	total := len(accounts)
	for i := 0; i < 4; i++ {
		go func(eId int) { // 一意なIDでゴルーチンを生成
			for j := 1; j < 1000; j++ { // ランダムに生成された送金を1000回実行
				from, to := rand.Intn(total), rand.Intn(total)
				for from == to {
					to = rand.Intn(total)
				}
				accounts[from].Transfer(&accounts[to], 10, eId) // 送金を実行
			}
			// 1000回の送金が完了したら、完了メッセージを表示
			fmt.Println(eId, "COMPLETE")
		}(i)
	}
	time.Sleep(60 * time.Second) // プログラムが終了する前に60秒待つ
}
