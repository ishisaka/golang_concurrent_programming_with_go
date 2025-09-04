// 条件成立後にチャネルをクローズする

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

// Broadcast は入力チャネルのデータを複数の出力チャネルへ均等に分配する関数です。
// 引数quitで終了判定を行い、指定された出力チャネル数nを生成してデータを転送します。
// 入力チャネルがクローズされた場合、生成されたすべての出力チャネルがクローズされます。
func Broadcast[k any](quit <-chan int, input <-chan k, n int) []chan k {
	// k型の出力チャネルをn個作成
	outputs := CreateAll[k](n)
	go func() {
		// 完了したらすべての出力チャネルをクローズ
		defer CloseAll(outputs...)
		var msg k
		mereData := true
		for mereData {
			select {
			case msg, mereData = <-input: // 入力チャネルからメッセージを読み込む
				// 入力チャネルがクローズされていなければメッセージを核出力チャネルへ書き込む
				if mereData {
					for _, output := range outputs {
						output <- msg
					}
				}
			case <-quit:
				return
			}
		}
	}()
	return outputs // 出力チャネルのスライスを返す
}

// CreateAll は指定された数のチャネルを生成し、それらをチャネルのスライスとして返します。
func CreateAll[k any](n int) []chan k {
	channels := make([]chan k, n)
	for i := range channels {
		channels[i] = make(chan k)
	}
	return channels
}

// CloseAll は複数のチャネルを一括でクローズする汎用関数です。
// 引数に指定されたすべてのチャネルを順番にクローズします。
// k はチャネルの要素型をジェネリックに指定可能です。
func CloseAll[k any](channels ...chan k) {
	for _, c := range channels {
		close(c)
	}
}

// frequentWords 関数は、単語チャネルから出現頻度の高い単語を集計し、ソートされた結果を出力するチャネルを返します。
// quitチャネルを利用して処理を中断することができます。
// 入力される単語を解析し、出現頻度の高い上位10個の単語を、カンマ区切りで出力チャネルに送信します。
func frequentWords(quit <-chan int, words <-chan string) <-chan string {
	mostFrequentWords := make(chan string)
	go func() {
		defer close(mostFrequentWords)
		// 個々の一意な単語の出現頻度を保存するマップを作成
		freqMap := make(map[string]int)
		// 一意な単語のリストを保存するスライスを作成
		freqList := make([]string, 0)
		metaData, word := true, ""
		for metaData {
			select {
			// 入力チャネルからの次のメッセージを読み込む
			case word, metaData = <-words:
				if metaData {
					// メッセージが新たな単語を含んでいれば、一意な単語をスライスへ追加
					if freqMap[word] == 0 {
						freqList = append(freqList, word)
					}
					freqMap[word] += 1 // 単語の出現頻度を1増やす
				}
			case <-quit:
				return
			}
			// すべての入力メッセージを読み込んだら、出現頻度順にリストをソート
			sort.Slice(freqList, func(a, b int) bool {
				return freqMap[freqList[a]] > freqMap[freqList[b]]
			})
		}
		// 出力チャネルへ出現頻度の高い10個の単語を書き込む
		mostFrequentWords <- strings.Join(freqList[:10], ", ")
	}()
	return mostFrequentWords
}

// Take は入力チャネルから指定された最大数の要素を取り出し、新しいチャネルに送信します。
// quit チャネルを使用して処理を中断できます。
func Take[K any](quit chan int, n int, input <-chan K) <-chan K {
	output := make(chan K)
	go func() {
		defer close(output)
		moreData := true
		var msg K
		// データがあり、カウントダウンが0よりも大きい限り、メッセージの転送を続ける
		for n > 0 && moreData {
			select {
			case msg, moreData = <-input: // 入力からメッセージを読み込む
				if moreData {
					output <- msg // メッセージを出力へ送信
					n--           // カウントダウン
				}
			case <-quit:
				return
			}
		}
		if n == 0 {
			close(quit)
		}
	}()
	return output
}

const downloaders = 20

func main() {
	quitWords := make(chan int) // Takeに使用するquitチャネルを作成、Take関数までは終了にこのチャネルを使う
	quit := make(chan int)      // quitチャネルを作成
	defer close(quit)
	urls := generateURLs(quitWords)
	pages := make([]<-chan string, downloaders)
	for i := range pages {
		pages[i] = downloadPages(quitWords, urls)
	}
	// 10000のカウントダウンでTake(n)ゴルーチンを作成して、extractWords の出力を読み込ませる
	words := Take(quitWords, 10000, extractWords(quitWords, FanIn(quitWords, pages...)))
	// パイプラインの残りに関してはquitチャネルを使う
	wordMulti := Broadcast(quit, words, 2)
	longestResults := longestWords(quit, wordMulti[0])
	frequentResults := frequentWords(quit, wordMulti[1])
	// 最も長い10個の単語を含む単一のメッセージを表示
	fmt.Println("Longest Words: ", <-longestResults)
	// もっとも出現頻度の高いメッセージを表示
	fmt.Println("Most Frequent Words: ", <-frequentResults)
}
