// charcountermutexslow.go
// Mutexによるロック範囲が適切では無い例
// ロック範囲が適切ではないためにゴルーチンの待ち時間が増え処理が全体として遅くなってしまう。
package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"
)

const AllLetters = "abcdefghijklmnopqrstuvwxyz"

func main() {
	mutex := sync.Mutex{} // Mutexの作成
	var frequency = make([]int, 26)
	for i := 1000; i <= 1030; i++ {
		url := fmt.Sprintf("https://rfc-editor.org/rfc/rfc%d.txt", i)
		go CountLetters(url, frequency, &mutex) // ゴルーチンにMutexのポインタを渡す
	}
	time.Sleep(60 * time.Second) // ゴルーチンの終了を待つ
	mutex.Lock()                 // Mutexで共有変数の読み込みを保護
	for i, c := range AllLetters {
		fmt.Printf("%c-%d ", c, frequency[i])
	}
	mutex.Unlock() // アンロック
}

// CountLetters
// 注：このプログラムでは、ロックとアンロックの配置が不適切なことを示すため、 意図的にmutexでgoroutine全体をロックして
// います。次のリストではこの問題を修正します。
func CountLetters(url string, frequency []int, mutex *sync.Mutex) {
	mutex.Lock() // 実行全体でMutexをロックして、全て逐次にする
	resp, _ := http.Get(url)
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		panic("Server returning error status code: " + resp.Status)
	}
	body, _ := io.ReadAll(resp.Body)
	for _, b := range body {
		c := strings.ToLower(string(b))
		cIndex := strings.Index(AllLetters, c)
		if cIndex >= 0 {
			frequency[cIndex] += 1
		}
	}
	fmt.Println("Completed:", url, time.Now().Format("15:04:05"))
	mutex.Unlock() // アンロック
}
