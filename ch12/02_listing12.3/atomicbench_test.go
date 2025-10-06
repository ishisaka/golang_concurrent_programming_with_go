package main

import (
	"sync/atomic"
	"testing"
)

var total = int64(0) // 64ビット整数を作成

func BenchmarkNormal(bench *testing.B) {
	for i := 0; i < bench.N; i++ {
		total += 1 // 通常の加算演算子を使ってtotal変数へ加算
	}
}

func BenchmarkAtomic(bench *testing.B) {
	for i := 0; i < bench.N; i++ {
		atomic.AddInt64(&total, 1) // アトミック加算操作関数を使ってtotal変数へ加算
	}
}
