package tree

import (
	"bytes"
	"fmt"
)

type Merger interface {
	Merge(interface{}, interface{}) interface{}
}

type SegmentTree struct {
	tree   []interface{}
	data   []interface{}
	merger Merger
}

// 线段树的初始化是自底向上进行的。
// 从每一个叶子结点开始（也就是原数组中的每一个元素），沿从叶子结点到根结点的路径向上按层构建。
// 在构建的每一步中，对应两个子结点的数据将被用来构建应当存储于它们母结点中的值。
// 每一个中间结点代表它的左右两个子结点对应区间融合过后的大区间所对应的值。
// 这个融合信息的过程可能依所需要处理的问题不同而不同
// （例如对于保存区间最小值的线段树来说，merge的过程应为min()函数）。
// 但从叶子节点（长度为1的区间）到根结点（代表输入的整个区间）更新的这一过程是统一的。

func NewSegmentTree(arr []interface{}, merger Merger) *SegmentTree {
	segmentTree := &SegmentTree{
		tree:   make([]interface{}, len(arr)*4),
		data:   make([]interface{}, len(arr)),
		merger: merger,
	}

	for i := 0; i < len(arr); i++ {
		segmentTree.data[i] = arr[i]
	}
	segmentTree.buildSegmentTree(0, 0, len(arr)-1)

	return segmentTree
}

// 在treeIndex的位置创建表示区间[l...r]的线段树
func (st *SegmentTree) buildSegmentTree(treeIndex int, l int, r int) {
	if l == r {
		st.tree[treeIndex] = st.data[l]
		return
	}
	leftTreeIndex := leftChild(treeIndex)
	rightTreeIndex := rightChild(treeIndex)

	mid := l + (r-l)/2
	st.buildSegmentTree(leftTreeIndex, l, mid)
	st.buildSegmentTree(rightTreeIndex, mid+1, r)

	st.tree[treeIndex] = st.merger.Merge(st.tree[leftTreeIndex], st.tree[rightTreeIndex])
}

func (st *SegmentTree) Size() int {
	return len(st.data)
}

func (st *SegmentTree) Get(index int) interface{} {
	if index < 0 || index >= len(st.data) {
		panic("Index is illegal.")
	}

	return st.data[index]
}

// 返回完全二叉树的数组表示中，一个索引所表示的元素的左孩子节点的索引
func leftChild(index int) int {
	return index*2 + 1
}

// 返回完全二叉树的数组表示中，一个索引所表示的元素的右孩子节点的索引
func rightChild(index int) int {
	return index*2 + 2
}

// 返回区间[queryL, queryR]的值
func (st *SegmentTree) Query(queryL int, queryR int) interface{} {
	if queryL < 0 || queryL >= len(st.data) || queryR < 0 || queryR > len(st.data) {
		panic("Index is illegal.")
	}

	return st.query(0, 0, len(st.data)-1, queryL, queryR)
}

// 在以treeIndex为根的线段树中[l...r]的范围里，搜索区间[queryL...queryR]的值
func (st *SegmentTree) query(treeIndex int, l int, r int, queryL int, queryR int) interface{} {
	if l == queryL && r == queryR {
		return st.tree[treeIndex]
	}

	mid := l + (r-l)/2
	// treeIndex的节点分为[l...mid]和[mid+1...r]两部分
	leftTreeIndex := leftChild(treeIndex)
	rightTreeIndex := rightChild(treeIndex)

	if queryL >= mid+1 {
		return st.query(rightTreeIndex, mid+1, r, queryL, queryR)
	} else if queryR <= mid {
		return st.query(leftTreeIndex, l, mid, queryL, queryR)
	}

	leftResult := st.query(leftTreeIndex, l, mid, queryL, mid)
	rightResult := st.query(rightTreeIndex, mid+1, r, mid+1, queryR)

	return st.merger.Merge(leftResult, rightResult)
}

func (st *SegmentTree) Set(index int, e interface{}) {
	if index < 0 || index > len(st.data) {
		panic("Index is Illegal")
	}

	st.data[index] = e
	st.set(0, 0, len(st.data)-1, index, e)
}

func (st *SegmentTree) set(treeIndex int, l int, r int, index int, e interface{}) {
	if l == r {
		st.tree[treeIndex] = e
		return
	}

	mid := l + (r-l)/2
	// treeIndex的节点分为[l...mid]和[mid+1...r]两部分
	leftTreeIndex := leftChild(treeIndex)
	rightTreeIndex := rightChild(treeIndex)
	if index >= mid+1 {
		st.set(rightTreeIndex, mid+1, r, index, e)
	} else {
		st.set(leftTreeIndex, l, mid, index, e)
	}

	st.tree[treeIndex] = st.merger.Merge(st.tree[leftTreeIndex], st.tree[rightTreeIndex])
}

func (st *SegmentTree) String() string {
	buffer := bytes.Buffer{}

	buffer.WriteString("[")
	for i := 0; i < len(st.tree); i++ {
		if st.tree[i] != nil {
			buffer.WriteString(fmt.Sprint(st.tree[i]))
		} else {
			buffer.WriteString("nil")
		}

		if i != len(st.tree)-1 {
			buffer.WriteString(", ")
		}
	}
	buffer.WriteString("]")

	return buffer.String()
}
