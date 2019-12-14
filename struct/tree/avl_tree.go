package tree

import (
	"bytes"
	"fmt"
	"math"

	"github.com/Zeb-D/go-util/struct/common"
)

type avlTreeNode struct {
	key         interface{}
	val         interface{}
	left, right *avlTreeNode
	height      int
}

// AVLTree 平衡二叉树
type AVLTree struct {
	root *avlTreeNode
	size int
}

// 生成 avlTreeNode 节点
func NewAVLTreeNode(k interface{}, v interface{}) *avlTreeNode {
	return &avlTreeNode{key: k, val: v, height: 1}
}

func NewAVLTree() *AVLTree {
	return &AVLTree{}
}

// 判断该二叉树是否是一颗二分搜索树
func (at *AVLTree) IsBST() bool {
	var keys []interface{}
	inOrderAVLTree(at.root, keys)

	for i := 1; i < len(keys); i++ {
		if ret, err := common.Compare(keys[i-1], keys[i]); err == nil && ret == 1 {
			return false
		}
	}
	return true
}

func inOrderAVLTree(n *avlTreeNode, keys []interface{}) {
	if n == nil {
		return
	}

	inOrderAVLTree(n.left, keys)
	keys = append(keys, n.key)
	inOrderAVLTree(n.right, keys)
}

// 判断该二叉树是否是一棵平衡二叉树
func (at *AVLTree) IsBalanced() bool {
	return at.isBalanced(at.root)
}

// 判断以Node为根的二叉树是否是一棵平衡二叉树，递归算法
func (at *AVLTree) isBalanced(n *avlTreeNode) bool {
	if n == nil {
		return true
	}

	balanceFactor := at.getBalanceFactor(n)
	if math.Abs(float64(balanceFactor)) > 1 {
		return false
	}
	return at.isBalanced(n.left) && at.isBalanced(n.right)
}

// 返回以 avlTreeNode 为根节点的二分搜索树中，key所在的节点
func (at *AVLTree) getNode(Node *avlTreeNode, key interface{}) *avlTreeNode {
	// 未找到等于 key 的节点
	if Node == nil {
		return nil
	}
	ret, err := common.Compare(key, Node.key)
	if err == nil && ret == 0 {
		return Node
	} else if err == nil && ret == -1 {
		return at.getNode(Node.left, key)
	} else {
		return at.getNode(Node.right, key)
	}
}

// 获得节点 avlTreeNode 的高度
func (at *AVLTree) getHeight(n *avlTreeNode) int {
	if n == nil {
		return 0
	}
	return n.height
}

// 获得节点 avlTreeNode 的平衡因子
func (at *AVLTree) getBalanceFactor(n *avlTreeNode) int {
	if n == nil {
		return 0
	}
	return at.getHeight(n.left) - at.getHeight(n.right)
}

// 对节点y进行向右旋转操作，返回旋转后新的根节点x
//        y                              x
//       / \                           /   \
//      x   T4     向右旋转 (y)        z     y
//     / \       - - - - - - - ->    / \   / \
//    z   T3                       T1  T2 T3 T4
//   / \
// T1   T2
func (at *AVLTree) rightRotate(y *avlTreeNode) *avlTreeNode {
	x := y.left
	T3 := x.right

	// 向右旋转过程
	x.right = y
	y.left = T3

	// 更新 height
	y.height = int(math.Max(float64(at.getHeight(y.left)), float64(at.getHeight(y.right)))) + 1
	x.height = int(math.Max(float64(at.getHeight(x.left)), float64(at.getHeight(x.right)))) + 1

	return x
}

// 对节点y进行向左旋转操作，返回旋转后新的根节点x
//    y                             x
//  /  \                          /   \
// T1   x      向左旋转 (y)       y     z
//     / \   - - - - - - - ->   / \   / \
//   T2  z                     T1 T2 T3 T4
//      / \
//     T3 T4
func (at *AVLTree) leftRotate(y *avlTreeNode) *avlTreeNode {
	x := y.right
	T2 := x.left

	// 向左旋转过程
	x.left = y
	y.right = T2

	// 更新 height
	y.height = int(math.Max(float64(at.getHeight(y.left)), float64(at.getHeight(y.right)))) + 1
	x.height = int(math.Max(float64(at.getHeight(x.left)), float64(at.getHeight(x.right)))) + 1

	return x
}

// 向二分搜索树中添加新的元素(key, value)
func (at *AVLTree) Add(key interface{}, val interface{}) {
	at.root = at.add(at.root, key, val)
}

// 向以node为根的二分搜索树中插入元素(key, value)，递归算法
// 返回插入新节点后二分搜索树的根
func (at *AVLTree) add(n *avlTreeNode, key interface{}, val interface{}) *avlTreeNode {
	if n == nil {
		at.size++
		return NewAVLTreeNode(key, val)
	}
	ret, err := common.Compare(key, n.key)
	if err == nil && ret < 0 {
		n.left = at.add(n.left, key, val)
	} else if err == nil && ret > 0 {
		n.right = at.add(n.right, key, val)
	} else {
		n.val = val
	}

	// 更新 height
	n.height = 1 + int(math.Max(float64(at.getHeight(n.left)), float64(at.getHeight(n.right))))
	// 计算平衡因子
	balanceFactor := at.getBalanceFactor(n)
	//if math.Abs(float64(balanceFactor)) > 1 {
	//	fmt.Println("unbalanced: ", balanceFactor)
	//}
	// 平衡维护
	// LL
	if balanceFactor > 1 && at.getBalanceFactor(n.left) >= 0 {
		return at.rightRotate(n)
	}
	// RR
	if balanceFactor < -1 && at.getBalanceFactor(n.right) <= 0 {
		return at.leftRotate(n)
	}
	// LR
	if balanceFactor > 1 && at.getBalanceFactor(n.left) < 0 {
		n.left = at.leftRotate(n.left)
		return at.rightRotate(n)
	}
	// RL
	if balanceFactor < -1 && at.getBalanceFactor(n.right) > 0 {
		n.right = at.rightRotate(n.right)
		return at.leftRotate(n)
	}
	return n
}

// 从二分搜索树中删除键为 key 的节点
func (at *AVLTree) Remove(key interface{}) interface{} {
	n := at.getNode(at.root, key)
	if n != nil {
		at.root = at.remove(at.root, key)
		return n.val
	}

	return nil
}

func (at *AVLTree) remove(n *avlTreeNode, key interface{}) *avlTreeNode {
	if n == nil {
		return nil
	}

	var retNode *avlTreeNode
	ret, err := common.Compare(key, n.key)
	if err == nil && ret < 0 {
		n.left = at.remove(n.left, key)
		retNode = n
	} else if err == nil && ret > 0 {
		n.right = at.remove(n.right, key)
		retNode = n
	} else {
		// 待删除节点左子树为空的情况
		if n.left == nil {
			rightNode := n.right
			n.right = nil
			at.size--
			retNode = rightNode
		} else
		// 待删除节点右子树为空的情况
		if n.right == nil {
			leftNode := n.left
			n.left = nil
			at.size--
			retNode = leftNode
		} else {
			// 待删除节点左右子树均不为空的情况

			// 找到比待删除节点大的最小节点, 即待删除节点右子树的最小节点
			// 用这个节点顶替待删除节点的位置
			successor := at.minimum(n.right)
			successor.right = at.remove(n.right, successor.key)
			successor.left = n.left

			n.left, n.right = nil, nil

			retNode = successor
		}
	}

	if retNode == nil {
		return nil
	}
	// 更新 height
	retNode.height = 1 + int(math.Max(float64(at.getHeight(retNode.left)), float64(at.getHeight(retNode.right))))
	// 计算平衡因子
	balanceFactor := at.getBalanceFactor(retNode)

	// 平衡维护
	// LL
	if balanceFactor > 1 && at.getBalanceFactor(retNode.left) >= 0 {
		return at.rightRotate(retNode)
	}
	// RR
	if balanceFactor < -1 && at.getBalanceFactor(retNode.right) <= 0 {
		return at.leftRotate(retNode)
	}
	// LR
	if balanceFactor > 1 && at.getBalanceFactor(retNode.left) < 0 {
		retNode.left = at.leftRotate(retNode.left)
		return at.rightRotate(retNode)
	}
	// RL
	if balanceFactor < -1 && at.getBalanceFactor(retNode.right) > 0 {
		retNode.right = at.rightRotate(retNode.right)
		return at.leftRotate(retNode)
	}
	return retNode
}

// 返回以node为根的二分搜索树的最小值所在的节点
func (at *AVLTree) minimum(n *avlTreeNode) *avlTreeNode {
	if n.left == nil {
		return n
	}
	return at.minimum(n.left)
}

func (at *AVLTree) Contains(key interface{}) bool {
	return at.getNode(at.root, key) != nil
}

func (at *AVLTree) Get(key interface{}) interface{} {
	n := at.getNode(at.root, key)
	if n == nil {
		return nil
	} else {
		return n.val
	}
}

func (at *AVLTree) Set(key interface{}, val interface{}) {
	n := at.getNode(at.root, key)
	if n == nil {
		panic(fmt.Sprintf("%v, doesn't exist", key))
	}

	n.val = val
}

func (at *AVLTree) Size() int {
	return at.size
}

func (at *AVLTree) IsEmpty() bool {
	return at.size == 0
}

func (at *AVLTree) String() string {
	var buffer bytes.Buffer
	avlTreeToString(at.root, 0, &buffer)
	return buffer.String()
}

// 生成以 avlTreeNode 为根节点，深度为 depth 的描述二叉树的字符串
func avlTreeToString(Node *avlTreeNode, depth int, buffer *bytes.Buffer) {
	if Node == nil {
		buffer.WriteString(generateDepthString(depth) + "nil\n")
		return
	}

	buffer.WriteString(generateDepthString(depth) + fmt.Sprintf("%s", Node.key) + "\n")
	avlTreeToString(Node.left, depth+1, buffer)
	avlTreeToString(Node.right, depth+1, buffer)
}

func generateDepthString(depth int) string {
	var buffer bytes.Buffer
	for i := 0; i < depth; i++ {
		buffer.WriteString("--")
	}
	return buffer.String()
}
