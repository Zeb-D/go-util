package heap

import (
	"errors"
	"fmt"
	"testing"
)

func TestPriorityQueue(t *testing.T) {
	k := 2
	nums := []int{1, 1, 1, 2, 2, 3}

	fmt.Println(topKFrequent(nums, k))

	k = 1
	nums = []int{1}
	fmt.Println(topKFrequent(nums, k))
}

/// 347. Top K Frequent Elements
/// https://leetcode.com/problems/top-k-frequent-elements/description/
///
/// 课程中在这里暂时没有介绍这个问题
/// 该代码主要用于使用Leetcode上的问题测试我们的MaxHeap类
func topKFrequent(nums []int, k int) []int {
	numMap := make(map[int]int)

	for _, num := range nums {
		numMap[num]++
	}

	pq := NewPriorityQueue(10)
	for num, f := range numMap {
		if pq.Size() < k {
			pq.Enqueue(&freq{e: num, frequency: f})
		} else {
			if f > pq.Front().(*freq).frequency {
				pq.Dequeue()
				pq.Enqueue(&freq{e: num, frequency: f})
			}
		}
	}

	var res []int
	for !pq.IsEmpty() {
		res = append(res, pq.Dequeue().(*freq).e)
	}
	return res
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
