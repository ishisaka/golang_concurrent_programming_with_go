// 練習問題4.1解答
package main

import (
	"fmt"
	"sync"
	"time"
)

func countdown(seconds *int, mutex *sync.Mutex) {
	mutex.Lock()
	sec := *seconds
	mutex.Unlock()
	for sec > 0 {
		time.Sleep(1 * time.Second)
		mutex.Lock()
		*seconds -= 1
		mutex.Unlock()
	}
}

func main() {
	mutex := sync.Mutex{}
	count := 5
	go countdown(&count, &mutex)
	c := count
	for c > 0 {
		time.Sleep(500 * time.Millisecond)
		mutex.Lock()
		fmt.Println(count)
		c = count
		mutex.Unlock()
	}
}
