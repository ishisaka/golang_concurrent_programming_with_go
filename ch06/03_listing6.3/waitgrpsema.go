// 5章の次作セマフォを使ってウェイトグループを自作した例

package listing6_3

import (
	"github.com/cutajarj/ConcurrentProgrammingWithGo/chapter5/listing5.16"
)

// WaitGrp はセマフォを利用して並行処理の待ち合わせを管理する構造体です。
type WaitGrp struct {
	// 5章で作成したセマフォへのポインタを格納
	sema *listing5_16.Semaphore
}

// NewWaitGrp は指定されたサイズでWaitGrp構造体を初期化し、そのポインタを返します。
// sizeはセマフォの許可数を管理するための値です。
func NewWaitGrp(size int) *WaitGrp {
	// 1 - sizeの許可数でセマフォを初期化
	return &WaitGrp{sema: listing5_16.NewSemaphore(1 - size)}
}

// Wait はセマフォを利用して全ての操作が完了するまで待機します。
func (wg *WaitGrp) Wait() {
	// セマフォに対してAcquire()を呼び出す
	wg.sema.Acquire()
}

// Done メソッドはセマフォに対して Release() を呼び出し、操作が完了したことを通知します。
func (wg *WaitGrp) Done() {
	// セマフォに対して Release()を呼び出す
	wg.sema.Release()
}
