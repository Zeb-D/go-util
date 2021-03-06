package hash

import (
	"fmt"
	"testing"
)

// q1:两数之和
func TestSolution1And2(t *testing.T) {
	nums := []int{2, 7, 11, 15}
	target := 9
	result := TwoSumSolution1(nums, target)
	fmt.Println(*result)

	result = TwoSumSolution2(nums, target)
	fmt.Println(*result)
}

//	暴力法 o(n^2)
func TwoSumSolution1(nums []int, target int) *[2]int {
	for i := 0; i < len(nums); i++ {
		for j := i + 1; j < len(nums); j++ {
			if nums[i]+nums[j] == target {
				result := [2]int{nums[i], nums[j]}
				return &result
			}
		}
	}

	return nil
}

//	一遍hash o(n)
func TwoSumSolution2(nums []int, target int) *[2]int {
	resultMap := make(map[int]int, len(nums))

	for i := 0; i < len(nums); i++ {
		if _, ok := resultMap[target-nums[i]]; ok {
			result := [2]int{resultMap[target-nums[i]], nums[i]}
			return &result
		}
		resultMap[nums[i]] = nums[i]
	}
	return nil
}
