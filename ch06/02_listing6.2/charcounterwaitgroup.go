// charcounterwaitgroup.go
// ウェイトグループの使用例
package main

import (
	"fmt"
	"sync"

	"github.com/cutajarj/ConcurrentProgrammingWithGo/chapter4/listing4.5"
)

func main() {
	wg := sync.WaitGroup{} // ウェイトグループの作成
	// 31の差分を加算
	wg.Add(31)
	mutex := sync.Mutex{}
	var frequency = make([]int, 26)
	for i := 1000; i <= 1030; i++ {
		url := fmt.Sprintf("https://rfc-editor.org/rfc/rfc%d.txt", i)
		// 無名関数でゴルーチンを作成
		go func() {
			listing4_5.CountLetters(url, frequency, &mutex)
			wg.Done() // 文字を数えおえたらDone()を呼び出す
		}()
	}
	wg.Wait() // 全てのゴルーチンの完了を待つ
	// ウェイトグループを使えば4章のコードのようにゴルーチンの終了を待つためのムダなスリープはいらない！
	mutex.Lock()
	for i, c := range listing4_5.AllLetters {
		fmt.Printf("%c-%d ", c, frequency[i])
	}
	mutex.Unlock()
}
