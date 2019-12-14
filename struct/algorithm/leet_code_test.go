package algorithm

import (
	"errors"
	"fmt"
	"github.com/Zeb-D/go-util/struct/heap"
	"testing"

	"github.com/Zeb-D/go-util/struct/Map"
	"github.com/Zeb-D/go-util/struct/set"
)

/// Leetcode 350. Intersection of Two Arrays II
/// https://leetcode.com/problems/intersection-of-two-arrays-ii/description/
func TestLeetCode350(t *testing.T) {
	num1 := []int{1, 2, 2, 1}
	num2 := []int{2, 2}
	fmt.Println(intersect(num1, num2))
}

// intersect 相同的元素
func intersect(nums1 []int, nums2 []int) []int {
	var res []int
	customMap := Map.NewBSTMap()

	for _, num := range nums1 {
		if customMap.Contains(num) {
			customMap.Set(num, customMap.Get(num).(int)+1)
		} else {
			customMap.Add(num, 1)
		}
	}

	for _, num := range nums2 {
		value := customMap.Get(num)

		if value != nil && value.(int) > 0 {
			res = append(res, num)
			customMap.Set(num, value.(int)-1)
		}
	}

	return res
}

/// Leetcode 349. Intersection of Two Arrays
/// https://leetcode.com/problems/intersection-of-two-arrays/description/
func TestLeetCode349(t *testing.T) {
	num1 := []int{1, 2, 2, 1}
	num2 := []int{2, 2}
	fmt.Println(intersection(num1, num2))
}

// intersection 相同的元素去重
func intersection(nums1 []int, nums2 []int) []int {
	var res []int
	bstSet := set.NewBSTSet()

	for _, num := range nums1 {
		bstSet.Add(num)
	}

	for _, num := range nums2 {
		if bstSet.Contains(num) {
			res = append(res, num)
			bstSet.Remove(num)
		}
	}

	return res
}

/// Leetcode 350. Intersection of Two Arrays II
/// https://leetcode.com/problems/intersection-of-two-arrays-ii/description/
func intersectMap(nums1 []int, nums2 []int) []int {
	var res []int
	nums1Map := make(map[int]int)

	for _, num := range nums1 {
		nums1Map[num]++
	}

	for _, num := range nums2 {
		if nums1Map[num] > 0 {
			res = append(res, num)
			nums1Map[num]--
		}
	}

	return res
}

func TestLeetCode350Map(t *testing.T) {
	num1 := []int{1, 2, 2, 1}
	num2 := []int{2, 2}
	fmt.Println(intersectMap(num1, num2))
}

/// 347. Top K Frequent Elements
/// https://leetcode.com/problems/top-k-frequent-elements/description/
///	查看哪个元素出现的次数最N多
/// 该代码主要用于使用Leetcode上的问题测试我们的MaxHeap类
func TopKFrequent(nums []int, k int) []int {
	numMap := make(map[int]int)

	for _, num := range nums {
		numMap[num]++
	}
	fmt.Println(numMap)
	maxHeap := heap.NewMaxHeap(len(numMap))
	for num, f := range numMap {
		//fmt.Println(num, "->", f)
		if maxHeap.Size() < k {
			maxHeap.Add(&freq{e: num, frequency: f})
		} else {
			if f > maxHeap.FindMin().(*freq).frequency {
				//maxHeap.ExtractMin()
				//maxHeap.Add(&freq{e: num, frequency: f})
				maxHeap.Replace(&freq{e: num, frequency: f})
			}
		}
	}

	var res []int
	for !maxHeap.IsEmpty() {
		res = append(res, maxHeap.ExtractMax().(*freq).e)
	}
	return res
}

func TestTopK(t *testing.T) {
	k := 2
	nums := []int{1, 1, 1, 2, 2, 3}

	fmt.Println("1->", TopKFrequent(nums, k))

	k = 1
	nums = []int{1, 2, 2}
	fmt.Println("2->", TopKFrequent(nums, k))
}

type freq struct {
	e         int
	frequency int
}

func (f *freq) Compare(b interface{}) (ret int, err error) {
	_, bok := b.(*freq)
	j := b.(*freq)
	//fmt.Println(bok, "<-aa->", j)
	if !bok { //类型不对
		err = errors.New("compare different type")
	}
	if f.frequency > j.frequency {
		ret = 1
	} else if f.frequency < j.frequency {
		ret = -1
	} else {
		ret = 0
	}

	return
}
