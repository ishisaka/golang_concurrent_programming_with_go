// 練習問題4.4の解答
package main

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"
)

func countLetters(url string, frequency map[string]int, mutex *sync.Mutex) {
	resp, _ := http.Get(url)
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		panic("Server's error: " + resp.Status)
	}
	body, _ := io.ReadAll(resp.Body)
	wordRegex := regexp.MustCompile(`[a-zA-Z]+`)
	mutex.Lock()
	for _, word := range wordRegex.FindAllString(string(body), -1) {
		wordLower := strings.ToLower(word)
		frequency[wordLower] += 1
	}
	mutex.Unlock()
	fmt.Println("Completed:", url)
}

func main() {
	mutex := sync.Mutex{}
	var frequency = make(map[string]int)
	for i := 1000; i <= 1020; i++ {
		url := fmt.Sprintf("https://rfc-editor.org/rfc/rfc%d.txt", i)
		go countLetters(url, frequency, &mutex)
	}
	time.Sleep(10 * time.Second)
	mutex.Lock()
	for k, v := range frequency {
		fmt.Println(k, "->", v)
	}
	mutex.Unlock()
}
