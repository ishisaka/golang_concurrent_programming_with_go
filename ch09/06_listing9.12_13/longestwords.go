// クローズ時に結果を出力する

package main

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"sort"
	"strings"
	"sync"
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

// FanIn は複数の入力チャネルを1つの出力チャネルに統合する関数です。
// quit チャネルで停止信号を受信すると動作を終了します。
// 入力されたすべてのチャネルがクローズされると、出力チャネルもクローズされます。
func FanIn[k any](quit <-chan int, allChannels ...<-chan k) chan k {
	wg := sync.WaitGroup{}   // ウェイトグループを作成
	wg.Add(len(allChannels)) // サイズを入力チャネルのか図と同じにする
	output := make(chan k)   // 出力チャネルを作成
	for _, c := range allChannels {
		go func(channel <-chan k) { // 入力チャネルごとにゴルーチンを作成
			defer wg.Done()
			for i := range channel {
				select {
				case output <- i: // 各受信メッセージを共有出力チャネルに送信
				case <-quit: // quitチャネルがクローズされたら、ゴルーチン終了
					return
				}
			}
		}(c) // 1つの入力チャネルをゴルーチンに渡す
	}
	go func() {
		wg.Wait()     // すべてのゴルーチンの終了待ち
		close(output) // それから出力チャネルをクローズ
	}()
	return output // 出力チャネルを返す
}

// longestWords 関数は単語チャネルから受け取った単語の中で最も長い10個の単語を選択し、出力チャネルに送信します。
// quitチャネルが提供されている場合、そのチャネルからの信号で処理を中断できます。
// 入力された単語は一意に管理され、最長の順でソートされます。
func longestWords(quit <-chan int, words <-chan string) <-chan string {
	longwords := make(chan string)
	go func() {
		defer close(longwords)
		// 一意な単語を保存するマップを作成
		uniqueWordsMap := make(map[string]bool)
		// あとでソートを簡単にするために一意な単語のリストを保存するスライスを作成
		uniqueWords := make([]string, 0)
		metaData, word := true, ""
		for metaData {
			select {
			case word, metaData = <-words:
				// チャネルがクローズされておらず新たな単語なら、マップとリストに追加
				if metaData && !uniqueWordsMap[word] {
					uniqueWordsMap[word] = true
					uniqueWords = append(uniqueWords, word)
				}
			case <-quit:
				return
			}
		}
		// 入力チャネルがクローズされたら、単語の長さで単語をソート
		sort.Slice(uniqueWords, func(i, j int) bool {
			return len(uniqueWords[i]) > len(uniqueWords[j])
		})
		// 出力チャネルへ10個の最も長い単語を含む文字列を送信
		longwords <- strings.Join(uniqueWords[:10], ", ")
	}()
	return longwords
}

const downloaders = 20

func main() {
	quit := make(chan int) // quitチャネルを作成
	defer close(quit)
	urls := generateURLs(quit)
	// ファンアウトパターンを使って処理を分散する
	// ダウンロードゴルーチンからの出力チャネルを保存するためのスライスを作成
	pages := make([]<-chan string, downloaders)
	for i := range pages {
		// ウェブページをダウンロードする20個のゴルーチンを作成し、出力チャネルを保存
		pages[i] = downloadPages(quit, urls)
	}
	// extractWordsの後にlongestWordsのゴルーチンをパイプラインに追加する
	result := longestWords(quit, extractWords(quit, FanIn(quit, pages...)))
	// 最も長い10個の単語を含む単一のメッセージを表示
	fmt.Println("Longest Words: ", <-result)
}
