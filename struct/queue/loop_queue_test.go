package queue

import (
	"fmt"
	"testing"
)

func TestLoopQueue_Dequeue(t *testing.T) {
	queue := NewLoopQueue(10)
	for i := 0; i < 10; i++ {
		queue.Enqueue(i)
		fmt.Println(queue)

		if i%3 == 2 {
			queue.Dequeue()
			fmt.Println(queue)
		}
	}
}

type TreeNode struct {
	Val   string
	Left  *TreeNode
	Right *TreeNode
}

type Pair struct {
	key   *TreeNode
	value int
}

// Leetcode 102. Binary Tree Level Order Traversal
// https://leetcode.com/problems/binary-tree-level-order-traversal/description/
// 二叉树的层序遍历
//
// 二叉树的层序遍历是一个典型的可以借助队列解决的问题。
// 该代码主要用于使用Leetcode上的问题测试我们的LoopQueue。
func TestPrintTreeLevelOrder(t *testing.T) {
	var res [][]string
	Node23f := &TreeNode{
		Val: "31-f",
	}
	Node2f := &TreeNode{
		Val:  "2-f",
		Left: Node23f,
	}
	Node2R := &TreeNode{
		Val: "2-r",
	}
	root := &TreeNode{
		Val:   "root",
		Left:  Node2f,
		Right: Node2R,
	}

	queue := NewLoopQueue(10)
	queue.Enqueue(Pair{root, 0})
	for !queue.IsEmpty() {
		front := queue.Dequeue().(Pair)
		node := front.key
		level := front.value

		if level == len(res) { // 每层第一次进行数组扩展
			res = append(res, []string{})
		}

		res[level] = append(res[level], node.Val)
		if node.Left != nil {
			queue.Enqueue(Pair{node.Left, level + 1})
		}
		if node.Right != nil {
			queue.Enqueue(Pair{node.Right, level + 1})
		}
	}

	fmt.Println(res)
}
