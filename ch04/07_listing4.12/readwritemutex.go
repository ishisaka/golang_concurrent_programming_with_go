// リーダー・ライター・ミューテックスの独自実装例
package listing4_12

import (
	"sync"
)

// Listing 4.12
type ReadWriteMutex struct {
	readersCounter int        // 現在クリティカルセクションにあるゴルーチンの数
	readersLock    sync.Mutex // リーダーのアクセスを同期するミューテックス
	globalLock     sync.Mutex // ライターのアクセスを待たせるミューテックス
}

// Listing 4.13
func (rw *ReadWriteMutex) ReadLock() {
	// 常に1つのゴルーチンしか許可されないようにアクセスを同期
	rw.readersLock.Lock()
	rw.readersCounter++ // ゴルーチンのカウンタを増加
	if rw.readersCounter == 1 {
		// リーダーゴルーチンが最初には行った場合、globalLockをロック
		rw.globalLock.Lock()
	}
	// 一度に1つのゴルーチンだけが許可されるようにアクセス同期していたのをアンロック
	rw.readersLock.Unlock()
}

func (rw *ReadWriteMutex) WriteLock() {
	rw.globalLock.Lock()
}

// Listing 4.14
func (rw *ReadWriteMutex) ReadUnlock() {
	// 常に1つのゴルーチンしか許可されないようにアクセスを同期
	rw.readersLock.Lock()
	rw.readersCounter-- // リーダーゴルーチンはreadersCounterの値を一つ減らす
	if rw.readersCounter == 0 {
		// リーダーゴルーチンが最後にでていった場合、globalLockをアンロック
		rw.globalLock.Unlock()
	}
	// 一度に1つのゴルーチンだけが許可されるようにアクセス同期していたのをアンロック
	rw.readersLock.Unlock()
}

func (rw *ReadWriteMutex) WriteUnlock() {
	rw.globalLock.Unlock()
}
