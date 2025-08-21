// gamesync.go
// Broadcast() の例。4人すべてのユーザーが接続するのを待ってから処理を進める。
package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	cond := sync.NewCond(&sync.Mutex{}) // 条件変数を作成
	playersInGame := 4                  // プレーヤーの数
	for playerId := 0; playerId < 4; playerId++ {
		// 条件変数とプレーヤーの数を共有するゴルーチンを起動
		go playerHandler(cond, &playersInGame, playerId)
		// 次のプレーヤーが起動するまでに1秒スリープ
		time.Sleep(1 * time.Second)
	}
}

func playerHandler(cond *sync.Cond, playersRemaining *int, playerId int) {
	cond.L.Lock() // 競合状態を避けるために条件変数のミューテックスをロック
	fmt.Println(playerId, ": Connected")
	*playersRemaining-- // 共有されたプレーヤーの数から1減らす
	if *playersRemaining == 0 {
		cond.Broadcast() // 全てのプーレーヤーが接続した場合、ブロードキャストを送信
	}
	for *playersRemaining > 0 {
		fmt.Println(playerId, ": Waiting for more players")
		cond.Wait() // 接続していない残りプレーヤーがいる限り待機
	}
	// ミューテックスをアンロックして全てのゴルーチンが実行を再開しゲーム開始
	cond.L.Unlock()
	fmt.Println("All players connected. Ready player", playerId)
	// ゲーム開始
}
