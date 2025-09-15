// ループキャリー依存性がある例

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
	files, _ := os.ReadDir(dir) // 指定されたディレクトリのファイル一覧を得る
	sha := sha256.New()         // ディレクト利用の新たな空のハッシュコンテナを作成
	for _, file := range files {
		if !file.IsDir() {
			fpath := filepath.Join(dir, file.Name())
			// ディレクトリ内の各ファイルのハッシュコードを計算
			hashOnFile := listing10_1.FHash(fpath)
			// 計算されたハッシュコードをディレクトリのハッシュコードへ結合
			sha.Write(hashOnFile)
		}
	}
	fmt.Printf("%s - %x\n", dir, sha.Sum(nil)) // ディレクトリのハッシュコードを表示
}
