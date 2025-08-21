// 条件変数を使ったリーダー・ライター・ロックの再検討
// スターベーションを回避するために条件変数を使用する例
// 読み込みを優先し、読み込み中もしくはその待機中のゴルーチンがある場合には書き込みは待たされる
package main

import (
	"fmt"
	"sync"
	"time"
)

// ReadWriteMutex はリーダー・ライターロックを管理するための構造体です。
// 複数のリーダーは同時にロックを保持できますが、ライターは単一スレッドのみが保持可能です。
// 条件変数を利用してリーダーとライター間のアクセス調整を行います。
type ReadWriteMutex struct {
	readersCounter int        // リーダーロックを現在保持しているリーダーの数
	writersWaiting int        // 現在待機しているライターの数
	writerActive   bool       // ライターがライターロックを保持しているかを示す
	cond           *sync.Cond // 条件変数
}

// NewReadWriteMutex は新しいReadWriteMutexインスタンスを作成して返します。
// このインスタンスはリーダー・ライターロックを管理するために使用されます。
func NewReadWriteMutex() *ReadWriteMutex {
	// 新たな条件変数と関連するミューテックスで新しいReadWriteMutexのインスタンスを作成
	return &ReadWriteMutex{cond: sync.NewCond(&sync.Mutex{})}
}

// ReadLock はリーダーアクセスを取得するためのメソッドです。ライターがアクティブまたは待機中の場合は待機します。
// リーダーロックを取得する際にリーダーのカウンタをインクリメントします。条件変数を使用して動作を管理します。
func (rw *ReadWriteMutex) ReadLock() {
	rw.cond.L.Lock() // ロック
	// ライターが待機中またはアクティブなら条件変数で待機
	for rw.writersWaiting > 0 || rw.writerActive {
		rw.cond.Wait()
	}
	rw.readersCounter++ // リーダーのカウンタを1増やす
	rw.cond.L.Unlock()  // アンロック
}

// WriteLock はライターロックを取得するメソッドです。他のリーダーやライターがロック中の場合は待機します。
// ライターの待機数をインクリメントし、リーダーや他のライターがアクティブでなくなるのを条件変数で待機します。
// ロック取得後、ライターのアクティブフラグを設定し、他スレッドがアクセスできないようにします。
func (rw *ReadWriteMutex) WriteLock() {
	rw.cond.L.Lock()    // ロック
	rw.writersWaiting++ // ライターの待機カウンタを1増やす
	// リーダーまたはライターがアクティブである限り条件変数で待機
	for rw.readersCounter > 0 || rw.writerActive {
		rw.cond.Wait()
	}
	rw.writersWaiting--    // 待機が終わったらライターの待機カウンタを1減らす
	rw.writerActive = true // 待機が終わったらライターのアクティブフラグを設定
	rw.cond.L.Unlock()     // アンロック
}

// ReadUnlock はリーダーがリーダーロックを解除する際に使用します。リーダーのカウンタをデクリメントし、最後のリーダーであれば条件変数をブロードキャストします。
func (rw *ReadWriteMutex) ReadUnlock() {
	rw.cond.L.Lock()    // ロック
	rw.readersCounter-- // リーダーのカウンタを1減らす
	// ゴルーチンが最後のリーダーであればブロードキャストを送信
	if rw.readersCounter == 0 {
		rw.cond.Broadcast()
	}
	rw.cond.L.Unlock() // アンロック
}

// WriteUnlock はライターがライターロックを解除する際に使用します。ブロードキャストを送信して他スレッドを通知します。
func (rw *ReadWriteMutex) WriteUnlock() {
	rw.cond.L.Lock()        // ロック
	rw.writerActive = false // ライターのアクティブフラグを解除
	rw.cond.Broadcast()     // ブロードキャストを送信
	rw.cond.L.Unlock()      // アンロック
}

func main() {
	rwMutex := NewReadWriteMutex()
	for i := 0; i < 2; i++ {
		go func() {
			for {
				rwMutex.ReadLock()
				time.Sleep(1 * time.Second)
				fmt.Println("Read done")
				rwMutex.ReadUnlock()
			}
		}()
	}
	time.Sleep(1 * time.Second)
	rwMutex.WriteLock()
	fmt.Println("Write finished")
}
