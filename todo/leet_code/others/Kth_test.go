package others

import (
	"fmt"
	"testing"
)

func TestTopK(t *testing.T) {
	//fmt.Println(topKFrequent([]int{1, 1, 1, 2, 2, 3}, 2))
	//fmt.Println(topKFrequent([]int{1}, 1))
	fmt.Println(topKFrequent([]int{-1, -1}, 1))
}

//	347. 前 K 个高频元素
//	输入: nums = [1,1,1,2,2,3], k = 2 输出: [1,2]
func topKFrequent(nums []int, k int) []int {
	freMap := make(map[int]int, k) //<元数、出现的次数>

	for _, v := range nums {
		val, _ := freMap[v]
		freMap[v] = val + 1
	}
	fmt.Println(freMap) //求出

	lengthNums := make([][]int, len(nums)+1) //对应出现的次数，作为下标，关联出现的元素
	for num, length := range freMap {
		lengthNums[length] = append(lengthNums[length], num)
	}
	fmt.Println(lengthNums) //次数越大，下标越往后，比如下标2为出现的频率，但是有多个元素出现了

	result := make([]int, 0)
	for end := len(lengthNums) - 1; end >= 0 && len(result) < k; end-- {
		value := lengthNums[end]
		if value == nil || len(value) == 0 {
			continue
		}
		restLength := k - len(result)
		if len(value) <= restLength {
			result = append(result, value...)
		} else {
			result = append(result, value[0:restLength]...)
		}
	}

	return result
}
