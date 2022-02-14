package todo

import (
	"fmt"
	"testing"
)

//专门复习算法练习

func TestQuickSort(t *testing.T) {
	arr := []int{3, 1, 0, 9, 3, 5, 8, 2, 4}
	fmt.Println(len(arr) - 1)
	quickSort(arr, 0, len(arr)-1)
	fmt.Println(arr)

	arr = []int{9, 8, 7, 6, 5, 4, 3, 2, 1}
	quickSort(arr, 0, len(arr)-1)
	fmt.Println(arr)
}

//快速排序
//选择一个关键值作为基准值。比基准值小的都在左边序列(一般是无序的)，
//比基准值大的都在右边(一般是无序的)。一般选择序列的第一个元素。
//一次循环:从后往前比较，用基准值和最后一个值比较，如果比基准值小的交换位置，如果没有 继续比较下一个，直到找到第一个比基准值小的值才交换。
//找到这个值之后，又从前往后开始比 较，如果有比基准值大的，交换位置，如果没有继续比较下一个，
//直到找到第一个比基准值大的 值才交换。直到从前往后的比较索引>从后往前比较的索引，结束第一次循环，此时，对于基准值 来说，左右两边就是有序的了。
func quickSort(arr []int, low, high int) {
	start, end := low, high
	key := arr[low] //基准值
	//先把基准值放到 左边小， 中间， 右边小
	for end > start {
		// 从后end往前比较
		for end > start && arr[end] >= key {
			end--
		}
		if arr[end] <= key {
			arr[start], arr[end] = arr[end], arr[start]
		}
		//从前往后比较
		for end > start && arr[start] <= key {
			start++
		}
		if arr[start] >= key {
			arr[start], arr[end] = arr[end], arr[start]
		}
	}

	// 再集中分别把 左边，右边，这部分单独再计算下
	if start > low { //说明产生了区间
		quickSort(arr, low, start-1)
	}
	if end < high {
		quickSort(arr, end+1, high)
	}
}

func TestInsertSort(t *testing.T) {
	arr := []int{3, 1, 0, 9, 3, 5, 8, 2, 4}
	insertSort(arr, true)
	fmt.Println(arr)
	insertSort(arr, false)
	fmt.Println(arr)
}

// 插入排序
//通过构建有序序列，对于未排序数据，在已排序序列中从后向前扫描，找到相应的位置并插入。
//插入排序非常类似于整扑克牌。在开始摸牌时，左手是空的，牌面朝下放在桌上。
func insertSort(arr []int, asc bool) []int {
	for i := 1; i < len(arr); i++ {
		insertVal := arr[i]
		index := i - 1
		if asc {
			for index >= 0 && insertVal < arr[index] {
				arr[index+1] = arr[index]
				index-- //让 index 向前移动
			}
		} else {
			for index >= 0 && insertVal > arr[index] {
				arr[index+1] = arr[index]
				index--
			}
		}
		arr[index+1] = insertVal //把插入的数放入合适位置
	}

	return arr
}

func TestBubbleSort(t *testing.T) {
	arr := []int{3, 1, 0, 9, 3, 5, 8, 2, 4}
	a := bubbleSort(arr, true)
	fmt.Println(arr)
	fmt.Println(a)
}

// 冒泡排序
//（1)比较前后相邻的二个数据，如果前面数据大于后面的数据，就将这二个数据交换。
//(2)这样对数组的第 0 个数据到 N-1 个数据进行一次遍历后，最大的一个数据就“沉”到数组第 N-1 个位置。
//(3)N=N-1，如果 N 不为 0 就重复前面二步，否则排序完成。
// asc 是否从小到大排序
func bubbleSort(arr []int, asc bool) []int {
	arrLength := len(arr)
	for i := 0; i < arrLength; i++ { //这是排序轮次
		for j := 1; j < arrLength-i; j++ { //这是排序数组下标
			//swap a[j] a[j-1]
			if asc {
				if arr[j-1] > arr[j] {
					arr[j-1], arr[j] = arr[j], arr[j-1]
				}
			} else {
				if arr[j-1] < arr[j] {
					arr[j-1], arr[j] = arr[j], arr[j-1]
				}
			}

		}
	}

	return arr
}

func TestSelectionSort(t *testing.T) {
	arr := []int{3, 1, 0, 9, 3, 5, 8, 2, 4}
	selectionSort(arr)
	fmt.Println(arr)
}

// 选择排序，与冒泡排序的区别：
// 原理一致，但交换方式不一样，先找到要排序的下标，然后最终进行换值，避免每次都换
func selectionSort(a []int) {
	n := len(a)
	for i := 0; i < n-1; i++ {
		selectIdx := i //放的最大？小值
		for j := i + 1; j < n; j++ {
			if a[j] < a[selectIdx] {
				selectIdx = j
			}
		}
		a[i], a[selectIdx] = a[selectIdx], a[i]
	}
}
