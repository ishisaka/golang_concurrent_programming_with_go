// パイプラインパターンの例

package main

import (
	"fmt"

	"github.com/cutajarj/ConcurrentProgrammingWithGo/chapter10/listing10.13"
)

// AddOnPipe はチャネルを通じてデータを処理する汎用的なパイプラインステップを作成します。
// q は終了シグナルとしてのチャネルで、受信時に処理を終了します。
// f は入力データを処理し、別の型に変換するための関数を指定します。
// in は処理対象データを受け取る入力チャネルを指定します。
// 処理結果は出力チャネルとして返されます。
func AddOnPipe[X, Y any](q <-chan int, f func(X) Y, in <-chan X) chan Y {
	output := make(chan Y) // Y型の出力チャネルを作成
	go func() {            // ゴルーチン開始
		defer close(output)
		for { // 無限ループ
			select { // selectを呼び出す
			case <-q: // quitチャネルがクローズされたらループを抜けてゴルーチンを収容
				return
			case input := <-in: // 利用可能なら入力チャネルからメッセージを受信
				output <- f(input) // f関数を呼び出し、関数の戻り値を出力チャネルへ出力
			}
		}
	}()
	return output
}

func main() {
	input := make(chan int) // 最初のステップに接続される最初の入力チャネルを作成
	quit := make(chan int)  // quitチャネルを作成
	// パイプラインの各ステップを接続し、各ステップの出力を次のステップの入力にする。
	output := AddOnPipe(quit, listing10_13.Box,
		AddOnPipe(quit, listing10_13.AddToppings,
			AddOnPipe(quit, listing10_13.Bake,
				AddOnPipe(quit, listing10_13.Mixture,
					AddOnPipe(quit, listing10_13.PrepareTray, input)))))
	go func() { // 10個のカップケーキを製造するためにパイプラインに対して10個の整数を送信するゴルーチンを作成
		for i := 0; i < 10; i++ {
			input <- i
		}
	}()
	for i := 0; i < 10; i++ { // パイプラインステップ最後の出力チャネルからメッセージ10箱分を読み取る
		fmt.Println(<-output, "received")
	}
}
