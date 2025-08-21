// semaaphore.go
// セマフォの実装例
package listing5_16

import (
	"sync"
)

type Semaphore struct {
	permits int        // セマフォに残っている許可数
	cond    *sync.Cond // 許可数が不足している場合に待機するために使う条件変数
}

func NewSemaphore(n int) *Semaphore {
	return &Semaphore{
		permits: n, // セマフォの初期許可数
		// 新たなセマフォの条件変数と関連するミューテックスの初期化
		cond: sync.NewCond(&sync.Mutex{}),
	}
}

func (rw *Semaphore) Acquire() {
	rw.cond.L.Lock() // permitsを保護するためにロック
	for rw.permits <= 0 {
		rw.cond.Wait() // 許可が使用可能になるまで待機
	}
	rw.permits--       // 使用可能な許可数を減らす
	rw.cond.L.Unlock() // アンロック
}

func (rw *Semaphore) Release() {
	rw.cond.L.Lock()   // permitsを保護するためにロック
	rw.permits++       // 使用可能な許可数を増やす
	rw.cond.Signal()   // 1つ以上の許可が利用可能であることを条件変数で通知
	rw.cond.L.Unlock() // アンロック
}
