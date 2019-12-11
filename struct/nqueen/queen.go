package nqueen

import (
	"bytes"
	"errors"
	"fmt"
	"math"
)

// IsLegalLoc 判断是否在同一行或同一列、是否在同斜线上
func IsLegalLoc(ql *QueenList, loc *Location) (bool, error) {
	if ql == nil {
		return false, errors.New("*QueenList is nil")
	}
	if loc == nil {
		return false, errors.New("*Location is nil")
	}
	//ql.RW.RLock()
	//defer ql.RW.Unlock()
	//遍历比较
	for head := ql.QList.Head(); head != nil; head = head.Next() {
		location := head.Value().(*Location)
		if location == nil {
			return true, errors.New("head.Value location is nil")
		}
		//判断是否在同一行或同一列
		if location.x == loc.x || loc.y == location.y {
			return false, errors.New(fmt.Sprintf("同一行或同一列[location->x:%d,->y:%d],[loc->x:%d,->y:%d]",
				location.x, location.y, loc.x, loc.y))
		} else //判断是否在同斜线上
		if math.Abs(float64(location.x-loc.x)) == math.Abs(float64(location.y-loc.y)) {
			return false, errors.New(fmt.Sprintf("同斜线[location->x:%d,->y:%d],[loc->x:%d,->y:%d]",
				location.x, location.y, loc.x, loc.y))
		}

	}

	return true, nil
}

// PrintQueenList 打印N皇后
func PrintQueenList(ql *QueenList, Size int) {
	if ql == nil {
		fmt.Println("ql is nil")
	}
	//ql.RW.RLock()
	//defer ql.RW.Unlock()
	var show = make([][]string, Size, Size)
	// init show array 0
	InitArray0(show, Size, Size)
	//fmt.Println(show)
	//fmt.Println(ql.QList)
	// fix array2 NQueen
	for head := ql.QList.Head(); head != nil && head.Value() != nil; head = head.Next() {
		loc := head.Value().(*Location)
		//fmt.Printf("[loc:%v]\n", loc)
		show[loc.x][loc.y] = "1"
	}
	// print
	PrintArray2(show, Size, Size)
	fmt.Println("---------------")
}

// InitArray0 二维数组 元素默认'0'
func InitArray0(array [][]string, x, y int) error {
	if array == nil {
		return errors.New("*[][]string is nil")
	}
	for i := 0; i < y; i++ {
		yArray := make([]string, x, x)
		for j := 0; j < x; j++ {
			yArray[j] = "0"
		}
		array[i] = yArray
	}
	return nil
}

// PrintArray2 输出二维数组，步骤，先遍历最外层数组，再遍历一唯数组
func PrintArray2(array [][]string, x, y int) {
	var buffer bytes.Buffer
	for i := 0; i < y; i++ {
		buffer.Reset()
		for j := 0; j < x; j++ {
			buffer.WriteString(array[i][j])
		}
		fmt.Println(buffer.String())
	}
}

// NQueen N皇后问题
func NQueen(ql *QueenList, startX, startY int, Size int) error {
	if ql == nil {
		return errors.New("ql == nil")
	}
	if ql.QList.Size() == Size { //当list元素个数为SIZE时，表示SIZE个皇后都摆放完毕，打印后即可退出函数。
		PrintQueenList(ql, Size) //打印皇后摆放方式
		return nil
	}
	for i := startX; i < Size; i++ {
		loc := NewLocation(i, startY)
		isLegal, _ := IsLegalLoc(ql, loc)
		//fmt.Printf("err:%v \n", err)
		if isLegal {
			ql.AddFirst(loc)              //将第y行的皇后摆放好
			NQueen(ql, 0, startY+1, Size) //开始摆放y+1行的皇后，同样从第0列开始摆放
			ql.RemoveFirst()              //每次摆放完一个皇后后，都要将其撤回，再试探其它的摆法。
		}
	}
	return nil
}
