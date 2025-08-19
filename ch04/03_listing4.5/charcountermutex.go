// charcountermutex.go
// 02_listing4.3_4のコードを修正し、Mutexによるロック範囲を適切にした例
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

func CountLetters(url string, frequency []int, mutex *sync.Mutex) {
	resp, _ := http.Get(url) // 関数の遅い部分はロックせずに並行に実行
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		panic("Server returning error code: " + resp.Status)
	}
	body, _ := io.ReadAll(resp.Body)
	mutex.Lock() // 関数の高速な部分だけをロックする
	for _, b := range body {
		c := strings.ToLower(string(b))
		cIndex := strings.Index(AllLetters, c)
		if cIndex >= 0 {
			frequency[cIndex] += 1
		}
	}
	mutex.Unlock()
	fmt.Println("Completed:", url, time.Now().Format("15:04:05"))
}

func main() {
	mutex := sync.Mutex{} // Mutexの作成
	var frequency = make([]int, 26)
	for i := 1000; i <= 1030; i++ {
		url := fmt.Sprintf("https://rfc-editor.org/rfc/rfc%d.txt", i)
		go CountLetters(url, frequency, &mutex) // ゴルーチンにMutexのポインタを渡す
	}
	time.Sleep(5 * time.Second) // ゴルーチンの終了を待つ。先ほどよりも大幅に待ち時間を減らせる。
	mutex.Lock()                // Mutexで共有変数の読み込みを保護
	for i, c := range AllLetters {
		fmt.Printf("%c-%d ", c, frequency[i])
	}
	mutex.Unlock() // アンロック
}
