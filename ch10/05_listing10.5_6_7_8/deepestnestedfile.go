package main

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type CodeDepth struct {
	file  string
	level int
}

// deepestNestedBlock は、指定されたファイル内のコードで最も深いネストレベルを計算し、結果をCodeDepth構造体で返します。
// ファイル名を指定する文字列filenameを受け取り、最も深いネストの整数値を持つCodeDepthを返します。
// ネストを表すレベルは波括弧('{', '}')に基づいて計算されます。
// ファイルを読み込む際にエラーが発生する場合は、そのエラーは無視されます。
// 閉じ波括弧に対応する開き波括弧が常に存在すると仮定しています。
func deepestNestedBlock(filename string) CodeDepth {
	code, _ := os.ReadFile(filename) // ファイル全体をバッファに読み込む
	max := 0
	level := 0
	for _, c := range code { // ファイル内のすべての単一文字を反復
		if c == '{' {
			level += 1 // 文字が開き波括弧ならlevelを1つ増やす
			// level変数の最大値を更新する
			max = int(math.Max(float64(max), float64(level)))
		} else if c == '}' {
			level -= 1 // 文字が閉じ波括弧ならlevelを1つ減らす
		}
	}
	return CodeDepth{filename, max} // ファイル名と最大ネストレベルを返す
}

// forkIfNeeded は、指定されたファイルがGoソースファイルである場合に新しいゴルーチンを生成して処理を分岐します。
// ファイルのパス、FileInfo、WaitGroup、結果を受け取るチャネルを引数として受け取ります。
// ファイルがディレクトリでないかつ拡張子が.goの場合にのみゴルーチンを実行します。
// deepestNestedBlock関数を呼び出し、その結果を共有チャネルに送信します。
// WaitGroupを使用して並行処理の完了を追跡および同期させます。
func forkIfNeeded(path string, info os.FileInfo,
	wg *sync.WaitGroup, results chan CodeDepth) {
	// パスがファイルで拡張子がgoのソースファイル化確認
	if !info.IsDir() && strings.HasSuffix(path, ".go") {
		wg.Add(1)   // ウェイトグループに1を加える
		go func() { // 新たなゴルーチンを作成
			// 関数を呼び出し、戻り値を共通の結果チャネルへ書き込む
			results <- deepestNestedBlock(path)
			wg.Done() // ウェイトグループで完了を通知する
		}()
	}
}

// joinResults は複数の部分結果を受け入れ、もっとも深いネストレベルを持つ結果を返すチャネルを生成します。
// 部分結果のチャネルがクローズされるまで結果を集約し、最終的な結果だけを出力します。
func joinResults(partialResults chan CodeDepth) chan CodeDepth {
	finalResult := make(chan CodeDepth) // 最終結果用のチャネルを作成
	max := CodeDepth{"", 0}
	go func() {
		// クローズされるまでチャネルからの結果を受信
		for pr := range partialResults {
			if pr.level > max.level { // もっとも深いネストブロックの値を記録
				max = pr
			}
		}
		finalResult <- max // チャネルのクローズ後、出力チャネルへ結果を書き込む
	}()
	return finalResult
}

func main() {
	dir := os.Args[1] // 引数からルートディレクトリを得る
	// フォークされたすべてのゴルーチンが使う共通チャネルを作成
	partialResults := make(chan CodeDepth)
	wg := sync.WaitGroup{}
	// ルートディレクトリを探索し、すべてのファイルに対して
	// ゴルーチンを作成するフォーク関数を呼び出す
	filepath.Walk(dir,
		func(path string, info os.FileInfo, err error) error {
			forkIfNeeded(path, info, &wg, partialResults)
			return nil
		})
	// ジョイン関数を呼び出し、最終結果を返すチャネルを得る
	finalResult := joinResults(partialResults)
	wg.Wait()             // フォークしたすべてのゴルーチンの作業完了を待つ
	close(partialResults) // 共通チャネルをクローズし、ジョインゴルーチンへ作業の完了を知らせる
	// 最終結果を出力チャネルから受信
	result := <-finalResult
	fmt.Printf("%s has the deepest nested code block of %d\n",
		result.file, result.level)
}
