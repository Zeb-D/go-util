package list

import (
	"fmt"
	"math"
	"testing"
)

/**
 * Definition for singly-linked list.
 * type ListNode struct {
 *     Val int
 *     Next *ListNode
 * }
 */
func swapPairs(head *ListNode) *ListNode {
	node := &ListNode{}
	node.Next = head
	//pre := node

	for head != nil && head.Next != nil {

		l1, l2 := head, head.Next
		next := l2.Next

		l1.Next = next
		l2.Next = l1
		//head = l2

		head = next

	}
	return node.Next
}

func initList(n int) *ListNode {
	root := &ListNode{}

	tmp := root
	for i := 1; i <= n; i++ {
		tmp.Next = &ListNode{Val: i}
		tmp = tmp.Next
	}
	return root
}

func TestLeetCode24(t *testing.T) {
	root := initList(4)
	fmt.Println(root)

	root = swapPairs(root)

	fmt.Println(root)
}

//指针方向替换
func reverseList(head *ListNode) *ListNode {
	tmp := &ListNode{}
	for head != nil {
		next := head.Next
		head.Next = tmp.Next
		tmp.Next = head

		head = next
	}
	return tmp.Next
}

func TestLeetCode206(t *testing.T) {
	root := initList(5)
	fmt.Println(root)
	root = reverseList(root)
	fmt.Println(root)
}

func isPalindrome(head *ListNode) bool {
	if head == nil || head.Next == nil {
		return true
	}
	//得到后一半的链表
	slow, fast := head, head.Next
	for fast != nil && fast.Next != nil {
		slow = slow.Next
		fast = fast.Next.Next
	}
	if fast != nil {
		slow = slow.Next
	}

	//获取到前一半的链表
	prev := head
	for head != slow {
		head = head.Next
	}

	//后一半链表倒置
	back := reverseList(slow)

	head.Next = nil

	//比较
	for prev != nil && back != nil {
		if prev.Val != back.Val {
			return false
		}
		prev = prev.Next
		back = back.Next
	}
	return true
}

func TestLeetCode234(t *testing.T) {
	root := initList(5)
	fmt.Println(root)
	fmt.Println(isPalindrome(root))
	math.Max(1+3, 3)
}

func rob(nums []int) int {
	if nums == nil || len(nums) == 0 {
		return 0
	}
	if len(nums) == 1 {
		return nums[0]
	}
	n := len(nums)
	return int(math.Max(rangeRob(nums, 0, n-2), rangeRob(nums, 1, n-1)))
}

func rangeRob(nums []int, start, end int) float64 {
	l1, l2 := float64(0), float64(0)
	for i := start; i <= end; i++ {
		l1, l2 = l2, math.Max(l1+float64(nums[i]), l2)
	}
	return l2
}
