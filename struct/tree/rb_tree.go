package tree

import (
	"bytes"
	"fmt"
	"github.com/Zeb-D/go-util/struct/common"
)

const RED = true
const BLACK = false

type rbTreeNode struct {
	key, value  interface{}
	left, right *rbTreeNode
	color       bool
}

// RBTree 红黑树
type RBTree struct {
	root *rbTreeNode
	size int
}

func NewRBTree() *RBTree {
	return &RBTree{}
}

// 生成node节点，默认为红色
func NewRBTreeNode(key interface{}, value interface{}) *rbTreeNode {
	return &rbTreeNode{key: key, value: value, color: RED}
}

// 返回以node为根节点的二分搜索树中，key所在的节点
func (rt *RBTree) getNode(n *rbTreeNode, key interface{}) *rbTreeNode {
	// 未找到等于 key 的节点
	if n == nil {
		return nil
	}

	ret, err := common.Compare(key, n.key)
	if err == nil && ret == 0 {
		return n
	} else if err == nil && ret < 0 {
		return rt.getNode(n.left, key)
	} else {
		return rt.getNode(n.right, key)
	}
}

// 判断节点node的颜色
func isRed(n *rbTreeNode) bool {
	if n == nil {
		return BLACK
	}
	return n.color
}

//   node                     x
//  /   \     左旋转         /  \
// T1   x   --------->   node   T3
//     / \              /   \
//    T2 T3            T1   T2
func leftRotate(n *rbTreeNode) *rbTreeNode {
	x := n.right

	// 左旋转
	n.right = x.left
	x.left = n

	x.color = n.color
	n.color = RED

	return x
}

//     node                   x
//    /   \     右旋转       /  \
//   x    T2   ------->   y   node
//  / \                       /  \
// y  T1                     T1  T2
func rightRotate(n *rbTreeNode) *rbTreeNode {
	x := n.left

	// 右旋转
	n.left = x.right
	x.right = n

	x.color = n.color
	n.color = RED

	return x
}

// 颜色翻转
func (node *rbTreeNode) flipColors() {
	node.color = RED
	node.left.color = BLACK
	node.right.color = BLACK
}

// 向红黑树中添加新的元素(key, value)
func (rt *RBTree) Add(key interface{}, val interface{}) {
	rt.root = rt.add(rt.root, key, val)
	rt.root.color = BLACK // 最终根节点为黑色节点
}

// 向以node为根的红黑树中插入元素(key, value)，递归算法
// 返回插入新节点后红黑树的根
func (rt *RBTree) add(n *rbTreeNode, key interface{}, val interface{}) *rbTreeNode {
	if n == nil {
		rt.size++
		return NewRBTreeNode(key, val)
	}
	ret, err := common.Compare(key, n.key)
	if err == nil && ret < 0 {
		n.left = rt.add(n.left, key, val)
	} else if err == nil && ret > 0 {
		n.right = rt.add(n.right, key, val)
	} else {
		n.value = val
	}

	if isRed(n.right) && !isRed(n.left) {
		n = leftRotate(n)
	}
	if isRed(n.left) && isRed(n.left.left) {
		n = rightRotate(n)
	}
	if isRed(n.left) && isRed(n.right) {
		n.flipColors()
	}

	return n
}

// 从二分搜索树中删除键为key的节点
func (rt *RBTree) Remove(key interface{}) interface{} {
	n := rt.getNode(rt.root, key)
	if n != nil {
		rt.root = rt.remove(rt.root, key)
		return n.value
	}

	return nil
}

func (rt *RBTree) remove(n *rbTreeNode, key interface{}) *rbTreeNode {
	if n == nil {
		return nil
	}
	ret, err := common.Compare(key, n.key)
	if err == nil && ret < 0 {
		n.left = rt.remove(n.left, key)
		return n
	} else if err == nil && ret > 0 {
		n.right = rt.remove(n.right, key)
		return n
	} else {
		// 待删除节点左子树为空的情况
		if n.left == nil {
			rightNode := n.right
			n.right = nil
			rt.size--
			return rightNode
		}
		// 待删除节点右子树为空的情况
		if n.right == nil {
			leftNode := n.left
			n.left = nil
			rt.size--
			return leftNode
		}
		// 待删除节点左右子树均不为空的情况

		// 找到比待删除节点大的最小节点, 即待删除节点右子树的最小节点
		// 用这个节点顶替待删除节点的位置
		successor := rt.minimum(n.right)
		successor.right = rt.removeMin(n.right)
		successor.left = n.left

		n.left, n.right = nil, nil

		return successor
	}
}

// 返回以node为根的二分搜索树的最小值所在的节点
func (rt *RBTree) minimum(n *rbTreeNode) *rbTreeNode {
	if n.left == nil {
		return n
	}
	return rt.minimum(n.left)
}

// 删除掉以node为根的二分搜索树中的最小节点
// 返回删除节点后新的二分搜索树的根
func (rt *RBTree) removeMin(n *rbTreeNode) *rbTreeNode {
	if n.left == nil {
		rightNode := n.right
		n.right = nil
		rt.size--

		return rightNode
	}

	n.left = rt.removeMin(n.left)
	return n
}

func (rt *RBTree) Contains(key interface{}) bool {
	return rt.getNode(rt.root, key) != nil
}

func (rt *RBTree) Get(key interface{}) interface{} {
	n := rt.getNode(rt.root, key)
	if n == nil {
		return nil
	} else {
		return n.value
	}
}

func (rt *RBTree) Set(key interface{}, val interface{}) {
	n := rt.getNode(rt.root, key)
	if n == nil {
		panic(fmt.Sprintf("%v, doesn't exist", key))
	}

	n.value = val
}

func (rt *RBTree) Size() int {
	return rt.size
}

func (rt *RBTree) IsEmpty() bool {
	return rt.size == 0
}

// 获得红黑树所有的 key
func (rt *RBTree) KeySet() []interface{} {
	var keySet []interface{}
	return recursive(rt.root, keySet)
}

func recursive(node *rbTreeNode, set []interface{}) []interface{} {
	if node == nil {
		return nil
	}

	recursive(node.left, set)
	recursive(node.right, set)
	return append(set, node.key)
}

func (rt *RBTree) String() string {
	var buffer bytes.Buffer
	rbTreeToString(rt.root, 0, &buffer)
	return buffer.String()
}

// 生成以node为根节点，深度为depth的描述二叉树的字符串
func rbTreeToString(node *rbTreeNode, depth int, buffer *bytes.Buffer) {
	if node == nil {
		buffer.WriteString(generateDepthString(depth) + "nil\n")
		return
	}

	buffer.WriteString(generateDepthString(depth) + fmt.Sprintf("%s", node.value) + "\n")
	rbTreeToString(node.left, depth+1, buffer)
	rbTreeToString(node.right, depth+1, buffer)
}
