// ループレベルへ移行パターンの例

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/cutajarj/ConcurrentProgrammingWithGo/chapter10/listing10.1"
)

func main() {
	dir := os.Args[1]
	files, _ := os.ReadDir(dir)
	wg := sync.WaitGroup{}
	for _, file := range files {
		if !file.IsDir() {
			wg.Add(1)
			// 反復でファイルのハッシュコードを計算するゴルーチンを開始
			go func(filename string) {
				fPath := filepath.Join(dir, filename)
				// List 10.1の関数を使用してファイルのファッ種コードを計算して表示する
				hash := listing10_1.FHash(fPath)
				fmt.Printf("%s: %x\n", filename, hash)
				wg.Done()
			}(file.Name())
		}
	}
	wg.Wait()
}
