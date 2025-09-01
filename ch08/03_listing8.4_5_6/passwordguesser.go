// defaultケースで並行計算を実行する

package main

import (
	"fmt"
	"time"
)

const (
	passwordToGuess = "go far"
	alphabet        = " abcdefghijklmnopqrstuvwxyz"
)

// toBase27は整数nを入力として受け取り、27進数表現の文字列に変換します。
// 各桁の値はalphabet変数に対応する文字で表されます。
// nが0より大きい値を持つ間、その値を27で割った余りを基に文字列を生成します。
// 最終的に生成される文字列は左側が上位桁、右側が下位桁となります。
func toBase27(n int) string {
	result := ""
	for n > 0 {
		result = string(alphabet[n%27]) + result
		n /= 27
	}
	return result
}

// guessPasswordは指定された範囲内でパスワードの推測を試みます。
// fromは推測を開始する範囲の最小値、uptoは範囲の最大値を指定します。
// stopチャネルを利用して他のゴルーチンの実行を停止することができます。
// パスワードが見つかった場合、resultチャネルに結果が送信されます。
// 通常の操作中にstopチャネルがクローズされた場合は、即座に処理を終了します。
func guessPassword(from int, upto int, stop chan int, result chan string) {
	// fromとuptoを始点と終点として、全てのパスワードの組合せをループ
	for guessN := from; guessN < upto; guessN++ {
		select {
		case <-stop:
			// stopチャネルでメッセージを受信すると、メッセージを出力して処理を停止
			fmt.Printf("Stopped at %d [%d, %d]\n", guessN, from, upto)
			return
		default:
			// パスワードの一致を確認
			if toBase27(guessN) == passwordToGuess {
				// 一致するパスワードをresultチャネルに送信
				result <- toBase27(guessN)
				// 他のゴルーチンが処理を中止するようにチャネルをクローズ
				close(stop)
				return
			}
		}
	}
	fmt.Printf("Not Found between %d and %d\n", from, upto)
}

func main() {
	// パスワードが発見されたことを通知する、大将ゴルーチンで使われる共通のチャネルを作成
	finished := make(chan int)
	// 発見されたパスワードを格納するチャネルを作成
	passwordFound := make(chan string)
	// 入力範囲[1, 1000万]. [1000万, 2000万]...を持つゴルーチンを作成
	for i := 0; i <= 387_420_488; i += 10_000_000 {
		go guessPassword(i, i+10_000_000, finished, passwordFound)
	}

	// パスワードが見つかるまで待機
	fmt.Println("password found: ", <-passwordFound)
	close(passwordFound)
	// 資源へのアクセスにパスワードを使うプログラムをシミュレート
	time.Sleep(5 * time.Second)
}
