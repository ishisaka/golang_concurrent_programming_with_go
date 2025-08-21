// 処理の途中でウェイトグループへの加算が必要な場合の例

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// fileSearch 関数は指定されたディレクトリ内のファイルを再帰的に検索し、一致したファイルパスを出力します。
// dir は検索を開始するディレクトリパスを指定します。
// filename は検索対象ファイル名の一部または全体を指定します。
// wg は並行処理の制御に使用される sync.WaitGroup を指します。
func fileSearch(dir string, filename string, wg *sync.WaitGroup) {
	// 関数に渡されたディレクトリ内のファイル全ての読み込み
	files, _ := os.ReadDir(dir)
	for _, file := range files {
		// 各ファイルをディレクトリに結合
		fpath := filepath.Join(dir, file.Name())
		if strings.Contains(file.Name(), filename) {
			// ファイル名が一致すればコンソールにパスを出力
			fmt.Println(fpath)
		}
		if file.IsDir() {
			// ディレクトリの場合、新たなゴルーチンを開始する前にウェイトグループに1を加算
			wg.Add(1)
			// ゴルーチンを再帰的に作成し、新たなディレクトリを検索
			go fileSearch(fpath, filename, wg)
		}
	}
	// 全てのファイルに対して処理が終わったら、ウェイトグループに対してDone()を呼ぶ
	wg.Done()
}

func main() {
	wg := sync.WaitGroup{} // 新たなウェイトグループを作成
	wg.Add(1)              // ウェイトグループに1を加算
	go fileSearch(os.Args[1], os.Args[2], &wg)
	wg.Wait() // 処理の完了を待機
}
