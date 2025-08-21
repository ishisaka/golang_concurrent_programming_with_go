// バリアの実装例

package listing6_10

import "sync"

// Barrier は複数のゴルーチン間で同期を取るためのバリア同期を提供する構造体です。
// 指定されたサイズ分のゴルーチンが到達するまで待機し、全て揃ったら進行を解除します。
// size フィールドはバリアに参加する対象ゴルーチンの総数を表します。
// waitCount フィールドは現在待機しているゴルーチン数を追跡するカウンタです。
// cond フィールドは条件変数を使って待機状態や通知を管理します。
type Barrier struct {
	size      int        // バリアへの参加対象の合計
	waitCount int        // 現在一時停止している実効の数を表すカウンタ変数｀
	cond      *sync.Cond // バリアで使われる条件変数
}

// NewBarrier は指定されたサイズのバリアを新たに作成し、そのポインタを返します。
// サイズはバリアに参加するゴルーチンの総数を指定します。
func NewBarrier(size int) *Barrier {
	condVar := sync.NewCond(&sync.Mutex{}) // 新たな条件変数の作成
	return &Barrier{size, 0, condVar}      // 新たなバリアを作成してポインタを返す
}

// Wait はバリア同期の主要なメソッドで、指定されたゴルーチン数が揃うまで現在のゴルーチンを待機させます。
// バリアサイズに到達すると、全ての待機中のゴルーチンを解除し進行を再開します。
func (b *Barrier) Wait() {
	b.cond.L.Lock()  // ミューテックスを使ってwaitCountへのアクセスを保護
	b.waitCount += 1 // カウントを1増やす
	// waitCountがバリアサイズに到達したらリセットして、条件変数に対してビロードキャストする
	if b.waitCount == b.size {
		b.waitCount = 0
		b.cond.Broadcast()

	} else {
		// waitCountがバリアサイズに到達していなければ待機
		b.cond.Wait()
	}
	b.cond.L.Unlock() // アンロック
}
