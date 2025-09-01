// nilチャネルパターン

package main

import (
	"fmt"
	"math/rand"
	"time"
)

// generateAmounts は指定された個数のランダムな金額を生成するチャネルを返します。
// n は生成する金額の数を指定します。
// 戻り値は生成された金額を受け取るための読み取り専用チャネルです。
func generateAmounts(n int) <-chan int {
	amounts := make(chan int) // 出力チャネルを作成
	go func() {
		defer close(amounts) // 終わったらチャネルを閉じる
		for i := 0; i < n; i++ {
			// 100ミリ秒ごとに出力チャネルに１から100までのランダムな値を書き込む
			amounts <- rand.Intn(100) + 1
			time.Sleep(100 * time.Millisecond)
		}
	}()
	return amounts
}

func main() {
	sales := generateAmounts(50)    // salesチャネルへ50個の売り上げを作成
	expenses := generateAmounts(40) // expensesチャネルへ40個の経費を生成
	endOfDayAmount := 0
	// nilではないチャネルがある間はループを続ける
	for sales != nil || expenses != nil {
		select {
		// salesチャネルから次の売り上げとチャネルのオープンフラグを読み込む
		case sale, moreData := <-sales:
			if moreData {
				fmt.Println("Sale of:", sale)
				endOfDayAmount += sale // 売り上げを1日の残高に加算
			} else {
				// チャネルがクローズしたらチャネルをnilにする
				sales = nil
			}
			// expensesチャネルから次の経費とチャネルのオープンフラグを読み込む
		case expense, moreData := <-expenses:
			if moreData {
				fmt.Println("Expense of:", expense)
				endOfDayAmount -= expense // 経費額を1日の残高から減算
			} else {
				// チャネルがクローズしたらチャネルをnilにする。
				expenses = nil
			}
		}
	}
	fmt.Println("End of day profit and loss:", endOfDayAmount)
}
