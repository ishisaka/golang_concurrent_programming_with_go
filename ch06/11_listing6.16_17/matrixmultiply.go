package main

import (
	"fmt"
	"math/rand"

	"github.com/cutajarj/ConcurrentProgrammingWithGo/chapter6/listing6.10"
)

const matrixSize = 3

func generateRandMatrix(matrix *[matrixSize][matrixSize]int) {
	for row := 0; row < matrixSize; row++ {
		for col := 0; col < matrixSize; col++ {
			matrix[row][col] = rand.Intn(10) - 5
		}
	}
}

func rowMultiply(matrixA, matrixB, result *[matrixSize][matrixSize]int,
	row int, barrier *listing6_10.Barrier) {
	for { // 無限ループ開始
		barrier.Wait() // main()ルーチンが行列を読み込むまでバリアで待機
		for col := 0; col < matrixSize; col++ {
			sum := 0
			for i := 0; i < matrixSize; i++ {
				// ゴルーチンの行の結果を計算
				sum += matrixA[row][i] * matrixB[i][col]
			}
			// ゴルーチンでの行の結果を正しい行と列に代入
			result[row][col] = sum
		}
		barrier.Wait() // 全ての計算が終わるまでバリアで待機
	}
}

func main() {
	var matrixA, matrixB, result [matrixSize][matrixSize]int
	// 行ごとのゴルーチンとmain()のゴルーチンを足した数をサイズとするバリアの作成
	barrier := listing6_10.NewBarrier(matrixSize + 1)
	for row := 0; row < matrixSize; row++ {
		// 行ごとのゴルーチンを作成して、正しい行番号を渡す
		go rowMultiply(&matrixA, &matrixB, &result, row, barrier)
	}

	for i := 0; i < 4; i++ {
		generateRandMatrix(&matrixA) // 両方の行列をランダムに作成s手読み込む
		generateRandMatrix(&matrixB)
		barrier.Wait() // ゴルーチンが計算できるようにバリアを開放
		barrier.Wait() // ゴルーチンが計算を完了するまで待機
		for i := 0; i < matrixSize; i++ {
			// 結果をコンソールへ出力
			fmt.Println(matrixA[i], matrixB[i], result[i])
		}
		fmt.Println()
	}
}
