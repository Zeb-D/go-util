package nqueen

import (
	"fmt"
	"testing"
)

func TestInitArray0(t *testing.T) {
	var XSize = 4
	var YSize = 4
	var show [][]string = make([][]string, YSize, YSize) //一唯数组个数
	fmt.Println(show)
	InitArray0(show, XSize, YSize)
	fmt.Println(show)
	show[2][2] = "1"
	PrintArray2(show, XSize, YSize)
}

func TestNQueen(t *testing.T) {
	var XSize int = 5
	fmt.Println("请输入N皇后解决方案")
	fmt.Scanln(&XSize)
	nq := NewQueenList()
	NQueen(nq, 0, 0, XSize)
	fmt.Printf("nq.QList:%v\n", nq.QList)
	//PrintQueenList(nq, 4)
}
