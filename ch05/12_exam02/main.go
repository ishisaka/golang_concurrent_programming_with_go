// 練習問題2の解答
package main

import (
	"fmt"
	"sync"
	"time"
)

// playerHandler は指定された条件変数を使用してプレイヤーの接続状況を管理します。
// cond は sync.Cond を示し、スレッド間の通信に使用されます。
// playersRemaining は接続待ちのプレイヤー数を表します。
// playerId は処理中のプレイヤーの識別子です。
// cancel はプレイヤー待機プロセスの中断フラグとして動作します。
func playerHandler(cond *sync.Cond, playersRemaining *int,
	playerId int, cancel *bool) {
	cond.L.Lock()
	fmt.Println(playerId, ": Connected")
	*playersRemaining--
	if *playersRemaining == 0 {
		cond.Broadcast()
	}
	for *playersRemaining > 0 && !*cancel {
		fmt.Println(playerId, ": Waiting for more players")
		cond.Wait()
	}
	cond.L.Unlock()
	if *cancel {
		fmt.Println(playerId, ": Game cancelled")
	} else {
		fmt.Println("All players connected. Ready player", playerId)
	}
}

// timeout は指定された sync.Cond を利用して、10秒後にキャンセルフラグを true に設定し、対応する条件変数に通知を送ります。
func timeout(cond *sync.Cond, cancel *bool) {
	time.Sleep(10 * time.Second)
	cond.L.Lock()
	*cancel = true
	cond.Broadcast()
	cond.L.Unlock()
}

func main() {
	cond := sync.NewCond(&sync.Mutex{})
	cancel := false
	go timeout(cond, &cancel) // タイムアウトのゴルーチンを起動
	playersInGame := 5
	for i := 0; i < 4; i++ {
		go playerHandler(cond, &playersInGame, i, &cancel)
		time.Sleep(3 * time.Second)
	}
	time.Sleep(20 * time.Second)
}
