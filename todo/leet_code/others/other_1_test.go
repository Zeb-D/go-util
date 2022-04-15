package others

import (
	"fmt"
	"math/rand"
	"sort"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestArrayNesting(t *testing.T) {
	//fmt.Println(arrayNesting([]int{5,4,0,3,1,6,2}))//4
	fmt.Println(arrayNesting([]int{1, 0, 3, 4, 2})) // 3

	fmt.Println(byte('a'))
	fmt.Println(byte('z'))
	fmt.Println(byte(122))
	fmt.Println(shiftingLetters("abc", []int{3, 5, 9}))

	fmt.Println(majorityElement([]int{1, 2, 5, 9, 5, 9, 5, 5, 5}))
	fmt.Println(majorityElement([]int{3, 2}))
	fmt.Println(majorityElement([]int{1, 2, 1, 1, 2, 2}))
}

//摩尔投票
//在集合中寻找可能存在的多数元素，这一元素在输入的序列重复出现并占到了序列元素的一半以上；
// 在第一遍遍历之后应该再进行一个遍历以统计第一次算法遍历的结果出现次数，确定其是否为众数；
// 如果一个序列中没有占到多数的元素，那么第一次的结果就可能是无效的随机元素。
//每次将两个不同的元素进行「抵消」，如果最后有元素剩余，则「可能」为元素个数大于总数一半的那个。
func majorityElement(nums []int) int {
	x, cnt := -1, 0
	for _, val := range nums {
		if cnt == 0 {
			x = val
			cnt = 1
		} else {
			if x == val {
				cnt += 1
			} else {
				cnt += -1
			}
		}

	}
	cnt = 0
	for _, val := range nums {
		if val == x {
			cnt++
		}
	}
	if cnt > len(nums)/2 {
		return x
	} else {
		return -1
	}
}

// 数组嵌套 leetcode 565
func arrayNesting(nums []int) int {
	max := 0
	for i := 0; i < len(nums); i++ {
		if nums[i] != -1 {
			start := nums[i]
			curMax := 0
			for nums[start] != -1 {
				curMax++
				//tmp := start
				//start = nums[start]
				//nums[tmp] = -1
				nums[start], start = -1, nums[start]
			}
			if curMax > max {
				max = curMax
			}
		}

	}
	return max
}

// 对应
func shiftingLetters(s string, shifts []int) string {
	endIdx := int(byte('z'))
	bs := []byte(s)
	shift := 0
	for i := len(bs) - 1; i >= 0; i-- {
		shift = (shift + shifts[i]) % 26
		b := int(bs[i]) + shift
		if b < endIdx {
			bs[i] = byte(b)
		} else {
			bs[i] = byte(b - 26)
		}
	}
	return string(bs)
}

//桶排序，根据出现次数进行排序
func TestFrequencySort(t *testing.T) {
	var s string = "tree"
	b := []byte(s)
	//计数
	c := map[byte]int{}
	for _, v := range b {
		c[v]++
	}

	//映射:次数：byte
	m := make([][]int, len(b)+1)
	for by, c := range c {
		m[c] = append(m[c], int(by))
	}

	//倒数输出
	ret := make([]byte, 0)
	for i := len(m) - 1; i >= 0; i-- {
		for j := 0; j < len(m[i]); j++ {
			for jj := 0; jj < i; jj++ {
				ret = append(ret, byte(m[i][j]))
			}
		}
	}
	fmt.Println(string(ret))
}

//根据制定颜色数字进行排序
func sortColors(a []int) {
	//扫描法，进行高低值互换
	start, end := 0, len(a)-1
	j := -1
	for start <= end {
		if a[start] == 0 {
			j++
			a[j], a[start] = a[start], a[j]
			start++
		} else if a[start] == 2 {
			a[start], a[end] = a[end], a[start]
			end--
		} else {
			start++
		}
	}

}

func TestEraseOverlapIntervals(t *testing.T) {
	//a := [][]int{{-3035,30075},{1937,6906},{11834,20971},{44578,45600},{28565,37578},{19621,34415},{32985,36313},{-8144,1080},{-15279,21851},{-27140,-14703},{-12098,16264},{-36057,-16287},{47939,48626},{-15129,-5773},{10508,46685},{-35323,-26257}}
	a := [][]int{{0, 1}, {3, 4}, {1, 2}}

	fmt.Println(eraseOverlapIntervals(a))
}

func eraseOverlapIntervals(intervals [][]int) int {
	n := len(intervals)
	if n == 0 {
		return 0
	}
	quickSortV2(intervals, 0, n-1)
	//sort.Slice(intervals, func(i, j int) bool { return intervals[i][1] < intervals[j][1] })
	end, cnt := intervals[0][1], 1
	for i := 1; i < n; i++ {
		if intervals[i][0] < end {
			continue
		}
		end = intervals[i][1]
		cnt++
	}
	return n - cnt
}

func quickSortV2(a [][]int, low, high int) {
	if low < high {
		idx := partitionV2(a, low, high)
		quickSortV2(a, idx+1, high)
		quickSortV2(a, low, idx-1)
	}
}

func partitionV2(a [][]int, low, high int) int {
	rIdx := low + rand.Intn(high-low)
	if rIdx < high {
		a[low], a[rIdx] = a[rIdx], a[low]
	}
	pivot := a[low][1]
	j := low
	for i := low + 1; i <= high; i++ {
		if a[i][1] < pivot {
			j++
			a[i], a[j] = a[j], a[i]
		}
	}
	a[low], a[j] = a[j], a[low]
	return j
}

// 一种计算深度的方式：2*lg N
func TestMaxDepth(t *testing.T) {
	for i := 0; i < 100; i++ {
		fmt.Println(maxDepth(i))
	}

}

func maxDepth(n int) int {
	var depth int
	for i := n; i > 0; i >>= 1 {
		depth++
	}
	return depth * 2
}

//给定一个非负整数的列表，安排它们形成最大的数字
//有一个盲点：前缀一样的，短的在前面，否则按照字符串比较
func TestMaxInt(t *testing.T) {
	tmp := []int{2, 32, 321} //2,32,321//9, 1, 2, 22, 11, 66
	//321322 323212
	fmt.Println("32" > "321")
	//排序
	sort.Slice(tmp, func(i, j int) bool {
		if strings.HasPrefix(strconv.Itoa(tmp[i]), strconv.Itoa(tmp[j])) {
			return false
		} else if strings.HasPrefix(strconv.Itoa(tmp[j]), strconv.Itoa(tmp[i])) {
			return false
		}
		ij := fmt.Sprint(tmp[i], tmp[j])
		ji := fmt.Sprint(tmp[j], tmp[i])
		return ij > ji
	})
	fmt.Println(tmp)
	var ret string
	for _, value := range tmp {
		ret = fmt.Sprint(ret, value)
	}
	fmt.Println(ret)
}

//每个协程都有一个编号，以创建协程的顺序轮流打印N次
//思路：每个协程都有一个队列，当自己的队列有值后，再通知到下一个队列
func TestPrintNo(t *testing.T) {
	gSize := 5
	forSize := 10

	//chan
	chanNum := make([]chan struct{}, gSize)
	noNum := make([]string, gSize)
	for i := 0; i < gSize; i++ {
		b := byte('A') + byte(i)
		noNum[i] = string(byte(b))
		chanNum[i] = make(chan struct{}, 1)
	}
	fmt.Println(noNum)

	wg := sync.WaitGroup{}
	wg.Add(gSize)
	for i := 0; i < gSize; i++ {
		go func(ii int) {
			defer func() {
				wg.Done()
			}()
			//持有当前的锁
			for j := 0; j < forSize; j++ {
				<-chanNum[ii]

				fmt.Print(noNum[ii])

				if ii == gSize-1 {
					chanNum[0] <- struct{}{}
				} else {
					chanNum[ii+1] <- struct{}{}
				}
			}
		}(i)
	}

	chanNum[0] <- struct{}{}

	wg.Wait()

}

//一个数组找收益最大的买入和卖出点（股票的最大利润）
//思路：动态规划，局部最优
func TestMaxStock(t *testing.T) {
	stocks := make([]int, 7)
	stocks[0] = rand.Intn(7)
	rand.Seed(time.Now().Unix())
	for i := 1; i < cap(stocks); i++ {
		stocks[i] = (stocks[i-1] + rand.Intn(i+1)*rand.Intn(stocks[i-1]+1)) / (rand.Intn(i+1) + 1)
	}
	fmt.Println(stocks)

	//双指针移动
	min, maxDif := stocks[0], 0
	for i := 1; i < len(stocks); i++ {
		if stocks[i]-min > maxDif {
			maxDif = stocks[i] - min
		}
		if stocks[i] < min {
			min = stocks[i]
		}
	}
	fmt.Println(maxDif, min)
}

//两个整数数组，数组中没有重复元素。
//第一个数组按顺序入栈，判断第二个数组是否是该栈的弹出顺序。
//思路：使用栈模式，压一个栈，进而去判断
func TestIsPopByOrder(t *testing.T) {
	n1 := []int{201, 202, 203, 204, 205}
	n2 := []int{201, 202, 205, 204, 203}

	var stack []int //临时栈
	for i := 0; i < len(n2); i++ {
		stack = append(stack, n2[i])
		for len(stack) != 0 && stack[len(stack)-1] == n1[len(n1)-1] {
			//两个一起出栈
			stack = stack[0 : len(stack)-1]
			n1 = n1[0 : len(n1)-1]
		}
	}

	fmt.Println(len(stack) == 0) //没元素 表示

}

func TestIsValid(t *testing.T) {
	n1 := []int{1, 2, 3, 4, 5}
	fmt.Println(n1[0:0]) //左臂右开
	fmt.Println(IsValid("{}(){}[]"))
	fmt.Println(IsValid("{[]()}"))
}

// go 的字符串实际是 byte 类型组成的切片
func IsValid(s string) bool {
	brackets := map[rune]rune{')': '(', ']': '[', '}': '{'}
	var stack []rune

	for _, char := range s {
		if char == '(' || char == '{' || char == '[' {
			// 括号左半部直接入栈
			stack = append(stack, char)
		} else if len(stack) > 0 && brackets[char] == stack[len(stack)-1] {
			// 括号右半部：栈中有数据，且此元素与栈尾元素相同，出栈
			stack = stack[:len(stack)-1]
		} else {
			return false
		}
	}

	// 循环结束，栈中还有数据返回 false
	return len(stack) == 0
}

//随机打乱有序数组
func TestShuffle(t *testing.T) {
	a := []int{9, 1, 2, 22, 11, 66}
	rand.Seed(time.Now().Unix())
	last := len(a) - 1
	for i := 0; i < len(a); i++ {
		pre := rand.Intn(i + 1)
		a[last], a[pre] = a[pre], a[last]
	}

	fmt.Println(a)
}

//原地移除字符串中空格
func TestTrimAllSpace(t *testing.T) {
	s := " aa ww w a c "
	fmt.Println(strings.Trim(s, " "))
	fmt.Println(strings.Replace(s, " ", "", -1))
}
