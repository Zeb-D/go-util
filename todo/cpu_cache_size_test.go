package todo

import (
	"sync"
	"testing"
)

const matrixLength = 6400

func createMatrix(size int) [][]int64 {
	matrix := make([][]int64, size)
	for i := 0; i < size; i++ {
		matrix[i] = make([]int64, size)
	}
	return matrix
}

func BenchmarkMatrixCombination(b *testing.B) {
	matrixA := createMatrix(matrixLength)
	matrixB := createMatrix(matrixLength)

	for n := 0; n < b.N; n++ {
		for i := 0; i < matrixLength; i++ {
			for j := 0; j < matrixLength; j++ {
				matrixA[i][j] = matrixA[i][j] + matrixB[i][j]
			}
		}
	}
}

func BenchmarkMatrixReversedCombination(b *testing.B) {
	matrixA := createMatrix(matrixLength)
	matrixB := createMatrix(matrixLength)

	for n := 0; n < b.N; n++ {
		for i := 0; i < matrixLength; i++ {
			for j := 0; j < matrixLength; j++ {
				matrixA[i][j] = matrixA[i][j] + matrixB[j][i]
			}
		}
	}
}

//通过加入小的遍历矩形块后，我们的整体遍历速度已经是最初版本的3倍了
func BenchmarkMatrixReversedCombinationPerBlock(b *testing.B) {
	matrixA := createMatrix(matrixLength)
	matrixB := createMatrix(matrixLength)
	blockSize := 8

	for n := 0; n < b.N; n++ {
		for i := 0; i < matrixLength; i += blockSize {
			for j := 0; j < matrixLength; j += blockSize {
				for ii := i; ii < i+blockSize; ii++ {
					for jj := j; jj < j+blockSize; jj++ {
						matrixA[ii][jj] = matrixA[ii][jj] + matrixB[jj][ii]
					}
				}
			}
		}
	}
}

// 这里M需要足够大，否则会存在goroutine 1已经执行完成，而goroutine 2还未启动的情况
const M = 1000000

type SimpleStruct struct {
	n int
}

//BenchmarkStructureFalseSharing-8            538       2245798 ns/op
func BenchmarkStructureFalseSharing(b *testing.B) {
	structA := SimpleStruct{}
	structB := SimpleStruct{}
	wg := sync.WaitGroup{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		wg.Add(2)
		go func() {
			for j := 0; j < M; j++ {
				structA.n += 1
			}
			wg.Done()
		}()
		go func() {
			for j := 0; j < M; j++ {
				structB.n += 1
			}
			wg.Done()
		}()
		wg.Wait()
	}
}

type PaddedStruct struct {
	n int
	_ CacheLinePad
}

type CacheLinePad struct {
	_ [CacheLinePadSize]byte
}

const CacheLinePadSize = 64

//然后，我们实例化这两个结构体，并继续在单独的goroutine中访问两个变量。
// 这里M需要足够大，否则会存在goroutine 1已经执行完成，而goroutine 2还未启动的情况
const MM = 1000000

// BenchmarkStructurePadding-8                 793       1506534 ns/op
func BenchmarkStructurePadding(b *testing.B) {
	structA := PaddedStruct{}
	structB := SimpleStruct{}
	wg := sync.WaitGroup{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		wg.Add(2)
		go func() {
			for j := 0; j < MM; j++ {
				structA.n += 1
			}
			wg.Done()
		}()
		go func() {
			for j := 0; j < MM; j++ {
				structB.n += 1
			}
			wg.Done()
		}()
		wg.Wait()
	}
}

//在上述的代码中，有两个6400*6400的初始化数组矩阵A和B，将A和B的元素进行相加，
// 第一种方式是对应行列坐标相加，即matrixA[i][j] = matrixA[i][j] + matrixB[i][j]，
// 第二种方式是对称行列坐标相加，即matrixA[i][j] = matrixA[i][j] + matrixB[j][i]。
// 那这两种不同的相加方式，会有什么样的结果呢？
