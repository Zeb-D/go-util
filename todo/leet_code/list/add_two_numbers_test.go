package list

import (
	"fmt"
	"github.com/Zeb-D/go-util/struct/list"
	"math"
	"testing"
)

// q2:两数相加
func TestSolutionAddTwoNumbers(t *testing.T) {
	l1 := list.NewLinkedList()
	l1.AddLast(2)
	l1.AddLast(4)
	l1.AddLast(3)
	l1.AddLast(6)

	l2 := list.NewLinkedList()
	l2.AddLast(5)
	l2.AddLast(6)
	l2.AddLast(4)

	_ = addTwoNumbers2(l1, l2)
	//fmt.Println(l3)

	ll1 := ListNode{Val: 2}
	ll1.Next = &ListNode{Val: 4}
	ll1.Next.Next = &ListNode{Val: 3}
	ll1.Next.Next.Next = &ListNode{Val: 6}

	ll2 := ListNode{Val: 5}
	ll2.Next = &ListNode{Val: 6}
	ll2.Next.Next = &ListNode{Val: 4}

	l4 := addTwoNumbers(&ll1, &ll2)
	fmt.Println(l4)
}

func addTwoNumbers(l1 *ListNode, l2 *ListNode) *ListNode {
	temp := &ListNode{}
	rs := temp
	for l1 != nil || l2 != nil {
		a := 0
		b := 0

		if l1 != nil {
			a = l1.Val
		}
		if l2 != nil {
			b = l2.Val
		}
		temp.Val = a + b

		if l1 != nil {
			l1 = l1.Next
		}
		if l2 != nil {
			l2 = l2.Next
		}
		if l2 != nil || l1 != nil {
			temp.Next = &ListNode{}
		}
		temp = temp.Next
	}

	temp = rs
	for temp != nil {
		if temp.Val > 9 {
			temp.Val = temp.Val - 10
			if temp.Next == nil {
				temp.Next = &ListNode{Val: 1}
			} else {
				temp.Next.Val = temp.Next.Val + 1
			}
		}
		temp = temp.Next
	}

	return rs
}

//	第一次遍历：两个链表对应每个节点分别取和，若含有空节点则空节点取0，产生一个新链表。
//	第二次遍历：对取完和的新链表遍历，判断当前的val是否大于等于10，大于或等于则其自身-10其next加1，若next为空则新建0节点。
func addTwoNumbers2(l1 *list.LinkedList, l2 *list.LinkedList) *list.LinkedList {

	rs := list.NewLinkedList()

	maxLen := math.Max(float64(l1.Size()), float64(l2.Size()))

	for i := 0; i < int(maxLen); i++ {
		a := 0
		b := 0
		if val, err := l1.Get(i); err == nil {
			a = val.(int)
		}
		if val, err := l2.Get(i); err == nil {
			b = val.(int)
		}
		rs.AddLast(a + b)
	}

	//对大于10的位置，-10 后面加1
	for i := 0; i < rs.Size(); i++ {
		if val, _ := rs.Get(i); val.(int) > 9 {
			rs.Set(i, val.(int)-10)

			val2, err2 := rs.Get(i + 1)
			if err2 != nil || val2 == nil {
				rs.Set(i+1, 1)
			} else {
				rs.Set(i+1, val2.(int)+1)
			}
		}
	}

	return rs
}

type ListNode struct {
	Val  int
	Next *ListNode
}
