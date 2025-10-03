package listing11_3_4

import (
	"fmt"
	"sync"
)

type BankAccount struct {
	id      string
	balance int
	mutex   sync.Mutex
}

// NewBankAccount は新しいBankAccountインスタンスを生成し、初期残高を100に設定します。
func NewBankAccount(id string) *BankAccount {
	return &BankAccount{
		id:      id,
		balance: 100,
		mutex:   sync.Mutex{},
	}
}

// Transfer は、指定された金額を送信元アカウントから送信先アカウントへ転送します。
func (src *BankAccount) Transfer(to *BankAccount, amount int, exId int) {
	fmt.Printf("%d Locking %s's account\n", exId, src.id)
	src.mutex.Lock() // 移動元口座に対するミューテックスをロック
	fmt.Printf("%d Locking %s's account\n", exId, to.id)
	to.mutex.Lock()       // 移動先口座に対するミューテックスをロック｀
	src.balance -= amount // 移動元から試験を減らす
	to.balance += amount  // 移動先に資金を追加する
	to.mutex.Unlock()     // 両方の口座のミューテックスをアンロック
	src.mutex.Unlock()
	fmt.Printf("%d Unlocked %s and %s\n", exId, src.id, to.id)
}
