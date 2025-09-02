// チャネルとゴルーチンによるパイプライン化

package main

import "fmt"

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

// main 関数はgenerateURLs関数を呼び出し、生成されたURLを標準出力に表示します。
func main() {
	quit := make(chan int) // quitチャネルを作成
	defer close(quit)
	result := generateURLs(quit) // 出力チャネルでURLを返すゴルーチン関数の呼び出し
	for url := range result {    // 出力チャネルから全てのメッセージを読み込む
		fmt.Println(url) // 結果を出力
	}
}
