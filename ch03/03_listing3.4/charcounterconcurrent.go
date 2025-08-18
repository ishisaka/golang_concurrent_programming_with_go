// charcounterconcurrent.go
// 02_listing3.2_3のコードをゴルーチンで並行化
package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const allLetters = "abcdefghijklmnopqrstuvwxyz"

func countLetters(url string, frequency []int) {
	resp, _ := http.Get(url)
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		panic("Server returning error status code: " + resp.Status)
	}
	body, _ := io.ReadAll(resp.Body)
	for _, b := range body {
		c := strings.ToLower(string(b))
		cIndex := strings.Index(allLetters, c)
		if cIndex >= 0 {
			frequency[cIndex] += 1
		}
	}
	fmt.Println("Completed:", url)
}

func main() {
	var frequency = make([]int, 26) // このスライスはヒープ上に作られる
	for i := 1000; i <= 1030; i++ {
		url := fmt.Sprintf("https://rfc-editor.org/rfc/rfc%d.txt", i)
		go countLetters(url, frequency) // ゴルーチンを使って並行に関数を呼び出す
	}
	time.Sleep(10 * time.Second) // ゴルーチンの終了を待つ
	for i, c := range allLetters {
		fmt.Printf("%c-%d ", c, frequency[i]) // 各文字を頻度とともに表示
	}
}

// 注意）
// このコードでは競合状態に対する考慮がなされていないので
// 出力は正しい結果にはならない
