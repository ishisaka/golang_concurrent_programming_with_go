// ループキャリーの依存性に対応した例
package main

import (
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"

	"github.com/cutajarj/ConcurrentProgrammingWithGo/chapter10/listing10.1"
)

func main() {
	dir := os.Args[1]
	files, _ := os.ReadDir(dir)
	sha := sha256.New()
	var prev, next chan int
	for _, file := range files {
		if !file.IsDir() {
			// ゴルーチン終了したことを通知する次のチャネルを作成
			next = make(chan int)
			go func(filename string, prev, next chan int) {
				// ファイルのハッシュコードを計算
				fpath := filepath.Join(dir, filename)
				hashOnFile := listing10_1.FHash(fpath)
				if prev != nil { // ゴルーチンの最初の反復でなければ前の反復が通知してくるまで待つ
					<-prev
				}
				sha.Write(hashOnFile) // ディレクトリの部分ハッシュコードを計算
				next <- 0             // 終了を次の反復に通知
			}(file.Name(), prev, next)
			// nextチャネルをprevへ代入。次のゴルーチンは現在の反復からの通知を待つ
			prev = next
		}
	}
	<-next // 結果を表示する前に最後の反復の終了を待つ
	fmt.Printf("%x\n", sha.Sum(nil))
}
