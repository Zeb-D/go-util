package others

import (
	"fmt"
	"math"
	"strings"
	"testing"
)

func TestTwoSum(t1 *testing.T) {
	fmt.Println(twoSum([]int{1, 2, 3, 4, 5}, 6)) //加法合
	fmt.Println(judgeSquareSum(999999999))       // 平方之合
	fmt.Println(validPalindrome("abca"))         //是否回文字符串
	fmt.Println(findLongestWord("abpcplea",
		[]string{"ale", "apple", "monkey", "plea", "abpcplaaa", "abpcllllll", "abccclllpppeeaaaa"}))
}

func findLongestWord(s string, dictionary []string) string {
	longthStr := ""
	for i := 0; i < len(dictionary); i++ {
		l1, l2 := len(longthStr), len(dictionary[i])
		if l1 > l2 || (l1 == l2 && longthStr < dictionary[i]) {
			continue
		}
		if isSubStr(s, dictionary[i]) {
			longthStr = dictionary[i]
		}
	}
	return longthStr
}

func isSubStr(s, dictionary string) bool {
	i, j := 0, 0
	b1, b2 := []byte(s), []byte(dictionary)
	for i < len(s) && j < len(dictionary) {
		if b1[i] == b2[j] {
			j++
		} else {
			i++
		}
	}
	return j == len(dictionary)
}

func validPalindrome(s string) bool {
	bs := []byte(s)
	i, j := 0, len(bs)-1
	for i < j {
		if bs[i] != bs[j] {
			ii := i + 1
			jj := j - 1
			return isPalindrome(bs, ii, j) || isPalindrome(bs, i, jj)
		}
		i++
		j--
	}
	return true
}

func isPalindrome(bs []byte, i, j int) bool {
	for i < j {
		if bs[i] != bs[j] {
			return false
		}
		i++
		j--
	}
	return true
}

func reverseVowels(s string) string {
	filters := []byte{'a', 'e', 'i', 'o', 'u'}
	bs := []byte(s)
	start, end := 0, len(bs)-1
	for i := 0; start <= end; i++ {
		if !strings.Contains("aeiouAEIOU", string(bs[i])) {
			if contains(filters, bs[end]) {
				bs[start], bs[end] = bs[end], bs[start]
			} else {
				end--
			}
		} else {
			start++
		}
	}
	return string(bs)
}
func contains(bs []byte, b byte) bool {
	for i := 0; i < len(bs); i++ {
		if bs[i] == b {
			return true
		}
	}
	return false
}

func judgeSquareSum(c int) bool {
	left, right := 0, int(math.Sqrt(float64(c)))
	for left <= right {
		sum := left*left + right*right
		if sum == c {
			return true
		} else if sum > c {
			right--
		} else {
			left++
		}
	}
	return false
}

func twoSum(numbers []int, target int) []int {
	length := len(numbers)
	left, right := 0, length-1
	for i := 0; i < length; i++ {
		if numbers[left]+numbers[right] == target {
			return []int{left, right}
		}
		if numbers[left]+numbers[right] > target {
			right--
		}
		if numbers[left]+numbers[right] < target {
			left++
		}
	}
	return nil
}
