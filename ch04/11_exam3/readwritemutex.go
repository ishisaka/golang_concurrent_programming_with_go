// 練習問題4.3の解答
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

func (rw *ReadWriteMutex) TryWriteLock() bool {
	return rw.globalLock.TryLock()
}

// 課題の解答
func (rw *ReadWriteMutex) TryReadLock() bool {
	if rw.readersLock.TryLock() {
		globalSuccess := true
		if rw.readersCounter == 0 {
			globalSuccess = rw.globalLock.TryLock()
		}
		if globalSuccess {
			rw.readersCounter++
		}
		rw.readersLock.Unlock()
		return globalSuccess
	} else {
		return false
	}
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
