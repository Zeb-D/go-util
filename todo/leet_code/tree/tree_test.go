package tree

import (
	"fmt"
	"testing"
)

//二叉排序树插入
type tree struct {
	v           int   //元素
	left, right *tree //左右节点
}

func insert(t *tree, k int) *tree {
	if t == nil {
		t = &tree{v: k}
		return t
	}
	if k <= t.v {
		t.left = insert(t.left, k)
	} else {
		t.right = insert(t.right, k)
	}
	return t
}

//前序遍历 中左右
func preOrder(t *tree, ret *[]int) {
	if t == nil {
		return
	}
	*ret = append(*ret, t.v)
	preOrder(t.left, ret)
	preOrder(t.right, ret)
}

//中序遍历 左中右
func inOrder(t *tree, ret *[]int) {
	if t == nil {
		return
	}
	inOrder(t.left, ret)
	*ret = append(*ret, t.v)
	inOrder(t.right, ret)
}

func subOrder(t *tree, ret *[]int) {
	if t == nil {
		return
	}
	inOrder(t.left, ret)
	inOrder(t.right, ret)
	*ret = append(*ret, t.v)
}

func TestTreeInsert(t *testing.T) {
	log.Infof("a:%s,%v", "aa", nil)
	root := insert(nil, 1)
	insert(root, 2)
	root = insert(root, 1)
	root = insert(root, 3)
	root = insert(root, 0)
	fmt.Printf("%+v \n", root)

	var ret []int
	preOrder(root, &ret)
	fmt.Println(ret)

	ret = []int{}
	inOrder(root, &ret)
	fmt.Println(ret)

	ret = []int{}
	subOrder(root, &ret)
	fmt.Println(ret)
}
