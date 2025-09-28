package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

func main() {
	const pagesToDownload = 30
	linesOnPage := make(chan int) // ページごとの行数を受け取るチャネル
	finalResult := make(chan int) // 最終結果を受け取るチャネル
	for i := 1000; i < 1000+pagesToDownload; i++ {
		go func(id int) { // ゴルーチンを作成しページをダウンロードする
			url := fmt.Sprintf("https://rfc-editor.org/rfc/rfc%d.txt", id)
			fmt.Println("Downloading", url)
			resp, _ := http.Get(url)
			if resp.StatusCode != 200 {
				panic("Server's error: " + resp.Status)
			}
			bodyBytes, _ := io.ReadAll(resp.Body)
			linesOnPage <- strings.Count(string(bodyBytes), "\n")
			resp.Body.Close()
		}(i)
	}
	go func() { // 各ページの集計結果を合計するためのゴルーチンを作成
		totalLines := 0
		for i := 0; i < pagesToDownload; i++ {
			totalLines += <-linesOnPage // 行数を合計
		}
		finalResult <- totalLines // 最終結果をチャネルに送信
	}()
	fmt.Println("Total lines:", <-finalResult)
}
