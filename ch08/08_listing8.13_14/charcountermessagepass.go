// 前々章でメモリ共有で実装されていたコードをチャネルの仕様に変更したもの

package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

const allLetters = "abcdefghijklmnopqrstuvwxyz"

// countLetters は指定された URL のテキストデータ内のアルファベットの頻度を非同期にカウントしチャネルで結果を返します。
// URL から取得するデータは小文字に変換され、英字（a-z）のみカウント対象になります。
// 結果は 26 個の整数値を含むスライスで、各インデックスがアルファベットの a から z に対応します。
// HTTP ステータスコードが 200 でない場合は panic を引き起こします。
// 非同期処理が完了した後にチャネルが閉じられます。
func countLetters(url string) <-chan []int {
	result := make(chan []int)
	go func() {
		defer close(result)
		frequency := make([]int, 26)
		resp, _ := http.Get(url)
		defer resp.Body.Close()
		if resp.StatusCode != 200 {
			panic("Server returning error code: " + resp.Status)
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
		result <- frequency
	}()
	return result
}

// main は指定された一連の RFC ドキュメント URL に対し、アルファベットの出現頻度を非同期で集計します。
// 各 URL の頻度データを集め、それを統合して A-Z の全アルファベットの合計出現回数を表示します。
// 非同期に処理を行うため、複数のチャネルを使用してデータを受け取ります。
func main() {
	results := make([]<-chan []int, 0)
	totalFrequencies := make([]int, 26)
	for i := 1000; i <= 1030; i++ {
		url := fmt.Sprintf("https://rfc-editor.org/rfc/rfc%d.txt", i)
		results = append(results, countLetters(url))
	}
	for _, c := range results {
		frequencyResult := <-c
		for i := 0; i < 26; i++ {
			totalFrequencies[i] += frequencyResult[i]
		}
	}
	for i, c := range allLetters {
		fmt.Printf("%c-%d ", c, totalFrequencies[i])
	}
}
