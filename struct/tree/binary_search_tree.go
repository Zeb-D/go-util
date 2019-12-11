package tree

import (
	"bytes"
	"fmt"
	"github.com/Zeb-D/go-util/struct/common"
	"github.com/Zeb-D/go-util/struct/list"
	"github.com/Zeb-D/go-util/struct/queue"
)

// 二分查找树
type BSTreeNode struct {
	e     interface{}
	left  *BSTreeNode
	right *BSTreeNode
}

func NewBSTreeNode(e interface{}) *BSTreeNode {
	return &BSTreeNode{e: e}
}

type BSTree struct {
	root *BSTreeNode
	size int
}

func NewBSTree() *BSTree {
	return &BSTree{}
}

func (tree *BSTree) Size() int {
	return tree.size
}

func (tree *BSTree) IsEmpty() bool {
	return tree.size == 0
}

// 向二分搜索树中添加新的元素 e
func (tree *BSTree) Add(e interface{}) {
	tree.root = tree.add(tree.root, e)
}

// 向以 BSTreeNode 为跟的二分搜索树中插入元素 e，递归算法
// 返回插入新节点后二分搜索树的根
func (tree *BSTree) add(n *BSTreeNode, e interface{}) *BSTreeNode {
	if n == nil {
		tree.size++
		return NewBSTreeNode(e)
	}

	// 递归调用
	ret, err := common.Compare(e, n.e)
	if err != nil {
		fmt.Printf("BSTreeNode:%v, add e:%v failed \n", n, e)
	}
	if err == nil && ret < 0 {
		n.left = tree.add(n.left, e)
	} else if err == nil && ret > 0 {
		n.right = tree.add(n.right, e)
	}
	return n
}

// 看二分搜索树中是否包含元素 e
func (tree *BSTree) Contains(e interface{}) bool {
	return contains(tree.root, e)
}

// 看以 BSTreeNode 为根的二分搜索树是否包含元素 e，递归算法
func contains(n *BSTreeNode, e interface{}) bool {
	if n == nil {
		return false
	}
	ret, err := common.Compare(e, n.e)
	if err != nil {
		fmt.Printf("BSTreeNode:%v, add e:%v failed \n", n, e)
	}
	if err == nil && ret == 0 {
		return true
	} else if err == nil && ret < 0 {
		return contains(n.left, e)
	} else {
		return contains(n.right, e)
	}
}

// 二分搜索树的前序遍历
func (tree *BSTree) PreOrder() {
	preOrder(tree.root)
}

// 前序遍历以 BSTreeNode 为根的二分搜索树，递归算法
func preOrder(n *BSTreeNode) {
	if n == nil {
		return
	}

	fmt.Println(n.e)
	preOrder(n.left)
	preOrder(n.right)
}

// 二分搜索树的非递归前序遍历
func (tree *BSTree) PreOrderNR() {
	// 使用之前我们自己实现的数组栈
	stack := list.NewArrayStack(20)
	stack.Push(tree.root)

	for !stack.IsEmpty() {
		cur := stack.Pop().(*BSTreeNode)
		fmt.Println(cur.e)

		if cur.right != nil {
			stack.Push(cur.right)
		}
		if cur.left != nil {
			stack.Push(cur.left)
		}
	}
}

// 二分搜索树的中序遍历
func (tree *BSTree) InOrder() {
	inOrder(tree.root)
}

// 中序遍历以 BSTreeNode 为根的二分搜索树，递归算法
func inOrder(n *BSTreeNode) {
	if n == nil {
		return
	}

	inOrder(n.left)
	fmt.Println(n.e)
	inOrder(n.right)
}

// 二分搜索树的后序遍历
func (tree *BSTree) PostOrder() {
	postOrder(tree.root)
}

// 后序遍历以 BSTreeNode 为根的二分搜索树，递归算法
func postOrder(n *BSTreeNode) {
	if n == nil {
		return
	}

	postOrder(n.left)
	postOrder(n.right)
	fmt.Println(n.e)
}

// 二分搜索树的层序遍历
func (tree *BSTree) LevelOrder() {
	// 使用我们之前实现的循环队列
	loopQueue := queue.NewLoopQueue(20)
	loopQueue.Enqueue(tree.root)
	for !loopQueue.IsEmpty() {
		cur := loopQueue.Dequeue().(*BSTreeNode)
		fmt.Println(cur.e)

		if cur.left != nil {
			loopQueue.Enqueue(cur.left)
		}
		if cur.right != nil {
			loopQueue.Enqueue(cur.right)
		}
	}
}

// 寻找二分搜索树的最小元素
func (tree *BSTree) Minimum() interface{} {
	if tree.size == 0 {
		panic("BSTree is empty!")
	}
	return minimum(tree.root).e
}

// 返回以 BSTreeNode 为根的二分搜索树的最小值所在的节点
func minimum(n *BSTreeNode) *BSTreeNode {
	if n.left == nil {
		return n
	}
	return minimum(n.left)
}

// 寻找二分搜索树的最大元素
func (tree *BSTree) Maximum() interface{} {
	if tree.size == 0 {
		panic("BSTree is empty!")
	}
	return maximum(tree.root).e
}

// 返回以 BSTreeNode 为根的二分搜索树的最大值所在的节点
func maximum(n *BSTreeNode) *BSTreeNode {
	if n.right == nil {
		return n
	}
	return maximum(n.right)
}

// 从二分搜索树中删除最小值所在的节点，返回最小值
func (tree *BSTree) RemoveMin() interface{} {
	// 获得最小值
	ret := tree.Minimum()
	tree.root = tree.removeMin(tree.root)
	return ret
}

// 删除以 BSTreeNode 为根的二分搜索树中的最小节点
// 返回删除节点后新的二分搜索树的根
func (tree *BSTree) removeMin(n *BSTreeNode) *BSTreeNode {
	if n.left == nil {
		rightNode := n.right
		tree.size--
		return rightNode
	}
	n.left = tree.removeMin(n.left)
	return n
}

// 从二分搜索树中删除最小值所在的节点，返回最小值
func (tree *BSTree) RemoveMax() interface{} {
	// 获得最小值
	ret := tree.Maximum()
	tree.root = tree.removeMax(tree.root)
	return ret
}

// 删除以 BSTreeNode 为根的二分搜索树中的最小节点
// 返回删除节点后新的二分搜索树的根
func (tree *BSTree) removeMax(n *BSTreeNode) *BSTreeNode {
	if n.right == nil {
		leftNode := n.left
		tree.size--
		return leftNode
	}
	n.right = tree.removeMax(n.right)
	return n
}

// 从二分搜索树中删除元素为 e 的节点
func (tree *BSTree) Remove(e interface{}) {
	tree.root = tree.remove(tree.root, e)
}

// 删除以 BSTreeNode 为根的二分搜索树中值为 e 的节点，递归算法
// 返回删除节点后新的二分搜索树的根
func (tree *BSTree) remove(n *BSTreeNode, e interface{}) *BSTreeNode {
	if n == nil {
		return nil
	}

	ret, err := common.Compare(e, n.e)
	if err != nil {
		fmt.Printf("BSTreeNode:%v, add e:%v failed \n", n, e)
	}
	if err == nil && ret < 0 {
		n.left = tree.remove(n.left, e)
		return n
	} else if err == nil && ret > 0 {
		n.right = tree.remove(n.right, e)
		return n
	} else {
		// 待删除节点左子树为空的情况
		if n.left == nil {
			rightNode := n.right
			n.right = nil
			tree.size--
			return rightNode
		}
		// 待删除节点右子树为空的情况
		if n.right == nil {
			leftNode := n.left
			n.left = nil
			tree.size--
			return leftNode
		}
		// 待删除节点左右子树均不为空的情况
		// 找到比待删除节点大的最小节点，即待删除节点右子树的最小节点
		// 用这个节点顶替待删除节点的位置
		successor := minimum(n.right)
		successor.right = tree.removeMin(n.right)
		successor.left = n.left
		n.left = nil
		n.right = nil

		return successor
	}
}

func (tree *BSTree) String() string {
	var buffer bytes.Buffer
	BSTreeToString(tree.root, 0, &buffer)
	return buffer.String()
}

// 生成以 BSTreeNode 为根节点，深度为 depth 的描述二叉树的字符串
func BSTreeToString(n *BSTreeNode, depth int, buffer *bytes.Buffer) {
	if n == nil {
		buffer.WriteString(newDepthString(depth) + "nil\n")
		return
	}

	buffer.WriteString(newDepthString(depth) + fmt.Sprint(n.e) + "\n")
	BSTreeToString(n.left, depth+1, buffer)
	BSTreeToString(n.right, depth+1, buffer)
}

func newDepthString(depth int) string {
	var buffer bytes.Buffer
	for i := 0; i < depth; i++ {
		buffer.WriteString("--")
	}
	return buffer.String()
}
