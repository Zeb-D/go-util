package set

import "github.com/Zeb-D/go-util/struct/tree"

// base on binary search tree
type BSTSet struct {
	tree *tree.BSTree
}

func NewBSTSet() *BSTSet {
	return &BSTSet{
		tree.NewBSTree(),
	}
}

func (b *BSTSet) Add(e interface{}) {
	b.tree.Add(e)
}

func (b *BSTSet) Remove(e interface{}) {
	b.tree.Remove(e)
}

func (b *BSTSet) Contains(e interface{}) bool {
	return b.tree.Contains(e)
}

func (b *BSTSet) Size() int {
	return b.tree.Size()
}

func (b *BSTSet) IsEmpty() bool {
	return b.tree.IsEmpty()
}

func (b *BSTSet) String() string {
	return b.tree.String()
}
