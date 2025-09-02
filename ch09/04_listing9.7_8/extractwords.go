// チャネルとゴルーチンによるパイプライン化
// さらにパイプラインに処理を追加

package main

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
)

// generateURLs 関数はquitチャネルを監視しながら、RFCドキュメントのURLを生成して送信するチャンネルを返します。
// 引数として受け取るquitチャネルで停止信号を受信した場合、URL生成を終了します。
// 生成されるURLはhttps://rfc-editor.org/rfc/rfc100.txtからhttps://rfc-editor.org/rfc/rfc130.txtまでです。
func generateURLs(quit <-chan int) <-chan string {
	urls := make(chan string)
	go func() {
		defer close(urls) // 完了したら、出力チャネルをクローズ
		for i := 100; i <= 130; i++ {
			url := fmt.Sprintf("https://rfc-editor.org/rfc/rfc%d.txt", i)
			select {
			case urls <- url: // 出力チャネルへ31個のURLを書き込む
			case <-quit:
				return
			}
		}
	}()
	return urls // 出力チャネルを返す
}

// downloadPages は指定されたURLのチャネルからページをダウンロードし、内容を出力チャネルに送信します。
// quitチャネルで停止信号を受信した場合、動作を終了します。
func downloadPages(quit <-chan int, urls <-chan string) <-chan string {
	pages := make(chan string) // ダウンロードされたウェブページを格納する出力チャネルを作成
	go func() {
		defer close(pages) // 終了したら出力チャネルをクローズ
		metaData, url := true, ""
		for metaData {
			// 入力チャネルがあればselectを継続
			select {
			// 新たなメッセージとデータの有無を示すフラグで変数を更新
			case url, metaData = <-urls:
				if metaData {
					// 新たなURLメッセージが届いたら、そのページをダウンロードして、
					// テキストをpagesチャネルへ送信
					resp, _ := http.Get(url)
					if resp.StatusCode != 200 {
						panic("Server returning error status code: " + resp.Status)
					}
					body, _ := io.ReadAll(resp.Body)
					pages <- string(body)
					_ = resp.Body.Close()
				}
			case <-quit: // quitチャネルへメッセージが届いたらゴルーチンを終了
				return
			}
		}
	}()
	return pages
}

// extractWords関数は、与えられた文字列チャネルから単語を抽出し、新たなチャネルで返します。
// quitチャネルを利用して処理を中断することができます。
// pagesチャネルから受け取った各文字列を正規表現で解析して単語を抜き出し、結果をwordsチャネルへ送信します。
func extractWords(quit <-chan int, pages <-chan string) <-chan string {
	words := make(chan string)
	go func() {
		defer close(words)
		// 単語を抽出するための正規表現を作成
		wordRegex := regexp.MustCompile("[a-zA-Z]+")
		metaData, pg := true, ""
		for metaData {
			select {
			// 新たなメッセージとデータの有無を示すフラグで変数を更新
			case pg, metaData = <-pages:
				if metaData {
					for _, word := range wordRegex.FindAllString(pg, -1) {
						// 新たなテキストページを受信すると、正規表現ですべての単語を
						// 抽出し、出力チャネルへ単語を送信
						words <- word
					}
				}
			case <-quit: // quitチャネルにメッセージがあれば、ゴルーチン終了
				return
			}
		}
	}()
	return words // 出力チャネルを返す
}

func main() {
	quit := make(chan int) // quitチャネルを作成
	defer close(quit)
	// 既存のパイプラインへ、ページをダウンロードする新たなゴルーチンを追加
	result := extractWords(quit, downloadPages(quit, generateURLs(quit)))
	for result := range result {
		fmt.Println(result) // 結果を出力
	}
}
