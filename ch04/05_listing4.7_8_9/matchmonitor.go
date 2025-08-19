// matchmonitor.go
// リーダー・ライダー・ミューテックスを使用する前のコード
package main

import (
	"fmt"
	"strconv"
	"sync"
	"time"
)

func matchRecorder(matchEvents *[]string, mutex *sync.Mutex) {
	for i := 0; ; i++ {
		mutex.Lock() // ミューテックスでmatchEventsへのアクセスを保護
		// 200ミリ秒ごとに試合イベントを含むモック文字列を追加
		*matchEvents = append(*matchEvents,
			"Match event "+strconv.Itoa(i))
		mutex.Unlock() // ミューテックスをアンロック
		time.Sleep(200 * time.Millisecond)
		fmt.Println("Appended match event")
	}
}

func clientHandler(mEvents *[]string, mutex *sync.Mutex, st time.Time) {
	for i := 0; i < 100; i++ {
		mutex.Lock() // ミューテックスで試合イベントのリストへのアクセスを保護
		// 試合イベントのスライス全体をコピーし、クライアントへの応答をシミュレート
		allEvents := copyAllEvents(mEvents)
		mutex.Unlock() // ミューテックスをアンロック

		timeTaken := time.Since(st) // 開始からの所要時間を計測
		// クライアントのサービスに要した時間をコンソールに出力
		fmt.Println(len(allEvents), "events copied in", timeTaken)
	}
}

func copyAllEvents(matchEvents *[]string) []string {
	allEvents := make([]string, 0, len(*matchEvents))
	for _, e := range *matchEvents {
		allEvents = append(allEvents, e)
	}
	return allEvents
}

func main() {
	mutex := sync.Mutex{} // 新たなミューテックスを初期化
	var matchEvents = make([]string, 0, 10000)
	// 進行中の試合をシミュレートするために、多くのスライスで事前にイベントを埋める
	for j := 0; j < 10000; j++ {
		matchEvents = append(matchEvents, "Match event")
	}
	// 試合記録のゴルーチンを開始
	go matchRecorder(&matchEvents, &mutex)
	start := time.Now() // クライアントハンドラのゴルーチンの開始時刻を記録
	for j := 0; j < 5000; j++ {
		// 大量のクライアントハンドラのゴルーチンを起動
		go clientHandler(&matchEvents, &mutex, start)
	}
	time.Sleep(100 * time.Second)
}
