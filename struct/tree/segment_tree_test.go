package tree

import (
	"fmt"
	"testing"
)

type addMeger struct {
}

func (m *addMeger) Merge(i, j interface{}) interface{} {
	return i.(int) + j.(int)
}

func TestSegmentTree(t *testing.T) {
	nums := []int{1, 3, 5}
	obj := convert(nums)
	fmt.Println(obj)
	fmt.Println(obj.Query(0, 2).(int))
	obj.Set(1, 2)
	fmt.Println(obj)
	fmt.Println(obj.Query(0, 2).(int))
}

func convert(nums []int) *SegmentTree {
	data := make([]interface{}, len(nums))
	for i := 0; i < len(nums); i++ {
		data[i] = nums[i]
	}
	return NewSegmentTree(data, &addMeger{})
}
