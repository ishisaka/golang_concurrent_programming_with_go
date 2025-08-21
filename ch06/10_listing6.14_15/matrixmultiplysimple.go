// ゴルーチンを使わない変更前

package main

import (
	"fmt"
	"math/rand"
)

const matrixSize = 3

func generateRandMatrix(matrix *[matrixSize][matrixSize]int) {
	for row := 0; row < matrixSize; row++ {
		for col := 0; col < matrixSize; col++ {
			// 各行、各列に-5から4までの間の乱数を割り当てる
			matrix[row][col] = rand.Intn(10) - 5
		}
	}
}

func matrixMultiply(matrixA, matrixB, result *[matrixSize][matrixSize]int) {
	for row := 0; row < matrixSize; row++ { // 全ての行を反復
		for col := 0; col < matrixSize; col++ { // 全ての列を反復
			sum := 0
			for i := 0; i < matrixSize; i++ {
				// A行の値とBの列の値を乗算したものの合計
				sum += matrixA[row][i] * matrixB[i][col]
			}
			result[row][col] = sum // result行列をsumで更新
		}
	}
}

func main() {
	var matrixA, matrixB, result [matrixSize][matrixSize]int
	for i := 0; i < 4; i++ {
		generateRandMatrix(&matrixA)
		generateRandMatrix(&matrixB)
		matrixMultiply(&matrixA, &matrixB, &result)
		for i := 0; i < matrixSize; i++ {
			fmt.Println(matrixA[i], matrixB[i], result[i])
		}
		fmt.Println()
	}
}
