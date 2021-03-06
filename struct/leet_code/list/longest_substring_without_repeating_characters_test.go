package list

import (
	"fmt"
	"testing"
)

// q3:无重复字符的最长子串
func TestSolutionLengthOfLongestSubstring(t *testing.T) {
	fmt.Println(lengthOfLongestSubstring("abcabcbb"))
	fmt.Println(lengthOfLongestSubstring("bbbbb"))
	fmt.Println(lengthOfLongestSubstring("pwwkew"))
}

//	Hash+双指针滑动窗口 o(n)
func lengthOfLongestSubstring(s string) int {
	bs := []byte(s)
	charMap := make(map[byte]int, len(bs)) //	保存每个字符最新的位置
	left := 0                              // 左指针初始位置
	length := 0                            //最新的最大子窜长度
	// 右指针初始位置
	for right := 0; right < len(bs); right++ {
		index, ok := charMap[bs[right]]
		charMap[bs[right]] = right
		if ok && index >= left { //存在重复的char 则更新左指针
			left = index + 1
		}
		if right-left+1 > length { //
			length = right - left + 1
		}
	}
	return length
}
