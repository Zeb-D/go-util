package Map

import (
	"bytes"
	"fmt"
	"github.com/Zeb-D/go-util/struct/common"
)

type Node struct {
	key   interface{}
	val   interface{}
	left  *Node
	right *Node
}

type BSTMap struct {
	root *Node
	size int
}

func NewBSTMap() *BSTMap {
	return &BSTMap{}
}

// 返回以node为根节点的二分搜索树中，key所在的节点
func (b *BSTMap) getNode(n *Node, key interface{}) *Node {
	// 未找到等于 key 的节点
	if n == nil {
		return nil
	}
	ret, err := common.Compare(key, n.key)
	if err == nil && ret == 0 {
		return n
	} else if err == nil && ret < 0 {
		return b.getNode(n.left, key)
	} else {
		return b.getNode(n.right, key)
	}
}

// 向二分搜索树中添加新的元素(key, value)
func (b *BSTMap) Add(key interface{}, val interface{}) {
	b.root = b.add(b.root, key, val)
}

// 向以node为根的二分搜索树中插入元素(key, value)，递归算法
// 返回插入新节点后二分搜索树的根
func (b *BSTMap) add(n *Node, key interface{}, val interface{}) *Node {
	if n == nil {
		b.size++
		return &Node{
			key: key,
			val: val,
		}
	}
	ret, err := common.Compare(key, n.key)
	if err == nil && ret < 0 {
		n.left = b.add(n.left, key, val)
	} else if err == nil && ret > 0 {
		n.right = b.add(n.right, key, val)
	} else {
		n.val = val
	}

	return n
}

// 从二分搜索树中删除键为key的节点
func (b *BSTMap) Remove(key interface{}) interface{} {
	n := b.getNode(b.root, key)
	if n != nil {
		b.root = b.remove(b.root, key)
		return n.val
	}

	return nil
}

func (b *BSTMap) remove(n *Node, key interface{}) *Node {
	if n == nil {
		return nil
	}
	ret, err := common.Compare(key, n.key)
	if err == nil && ret < 0 {
		n.left = b.remove(n.left, key)
		return n
	} else if err == nil && ret > 0 {
		n.right = b.remove(n.right, key)
		return n
	} else {
		// 待删除节点左子树为空的情况
		if n.left == nil {
			rightNode := n.right
			n.right = nil
			b.size--
			return rightNode
		}
		// 待删除节点右子树为空的情况
		if n.right == nil {
			leftNode := n.left
			n.left = nil
			b.size--
			return leftNode
		}
		// 待删除节点左右子树均不为空的情况

		// 找到比待删除节点大的最小节点, 即待删除节点右子树的最小节点
		// 用这个节点顶替待删除节点的位置
		successor := b.minimum(n.right)
		successor.right = b.removeMin(n.right)
		successor.left = n.left

		n.left, n.right = nil, nil

		return successor
	}
}

// 返回以node为根的二分搜索树的最小值所在的节点
func (b *BSTMap) minimum(n *Node) *Node {
	if n.left == nil {
		return n
	}
	return b.minimum(n.left)
}

// 删除掉以node为根的二分搜索树中的最小节点
// 返回删除节点后新的二分搜索树的根
func (b *BSTMap) removeMin(n *Node) *Node {
	if n.left == nil {
		rightNode := n.right
		n.right = nil
		b.size--

		return rightNode
	}

	n.left = b.removeMin(n.left)
	return n
}

func (b *BSTMap) Contains(key interface{}) bool {
	return b.getNode(b.root, key) != nil
}

func (b *BSTMap) Get(key interface{}) interface{} {
	n := b.getNode(b.root, key)
	if n == nil {
		return nil
	} else {
		return n.val
	}
}

func (b *BSTMap) Set(key interface{}, val interface{}) {
	n := b.getNode(b.root, key)
	if n == nil {
		panic(fmt.Sprintf("%v, doesn't exist", key))
	}

	n.val = val
}

func (b *BSTMap) Size() int {
	return b.size
}

func (b *BSTMap) IsEmpty() bool {
	return b.size == 0
}

func (b *BSTMap) String() string {
	var buffer bytes.Buffer
	generateBSTString(b.root, 0, &buffer)
	return buffer.String()
}

// 生成以 Node 为根节点，深度为 depth 的描述二叉树的字符串
func generateBSTString(Node *Node, depth int, buffer *bytes.Buffer) {
	if Node == nil {
		buffer.WriteString(generateDepthString(depth) + "nil\n")
		return
	}

	buffer.WriteString(generateDepthString(depth) + fmt.Sprint(Node.val) + "\n")
	generateBSTString(Node.left, depth+1, buffer)
	generateBSTString(Node.right, depth+1, buffer)
}

func generateDepthString(depth int) string {
	var buffer bytes.Buffer
	for i := 0; i < depth; i++ {
		buffer.WriteString("--")
	}
	return buffer.String()
}
