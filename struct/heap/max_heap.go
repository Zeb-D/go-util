package heap

import (
	"github.com/Zeb-D/go-util/struct/common"
	"github.com/Zeb-D/go-util/struct/list"
)

type maxHeap struct {
	data *list.ArrayList
}

func (h *maxHeap) String() string {
	return h.data.String()
}

func NewMaxHeap(capacity int) *maxHeap {
	return &maxHeap{
		data: list.NewArrayList(capacity),
	}
}

// 返回堆中的元素个数
func (h *maxHeap) Size() int {
	return h.data.Size()
}

// 返回一个布尔值, 表示堆中是否为空
func (h *maxHeap) IsEmpty() bool {
	return h.data.IsEmpty()
}

// 返回完全二叉树的数组表示中，一个索引所表示的元素的父亲节点的索引
func parent(index int) int {
	if index == 0 {
		panic("index-0 doesn't have parent.")
	}
	return (index - 1) / 2
}

// 返回完全二叉树的数组表示中，一个索引所表示的元素的左孩子节点的索引
func leftChild(index int) int {
	return index*2 + 1
}

// 返回完全二叉树的数组表示中，一个索引所表示的元素的右孩子节点的索引
func rightChild(index int) int {
	return index*2 + 2
}

// 向堆中添加元素
func (h *maxHeap) Add(e interface{}) {
	h.data.AddLast(e)
	h.siftUp(h.data.Size() - 1)
}

func (h *maxHeap) siftUp(k int) (err error) {
	for k > 0 {
		ret, err := common.Compare(h.data.Get(k), h.data.Get(parent(k)))
		if err != nil {
			return err
		}
		if err == nil && ret > 0 {
			h.data.Swap(k, parent(k))
			k = parent(k)
		} else {
			return err
		}
	}
	return
}

func (h *maxHeap) FindMin() interface{} {
	if h.data.Size() == 0 {
		panic("Can not findMax when heap is empty.")
	}
	return h.data.Get(h.data.Size() - 1)
}

func (h *maxHeap) FindMax() interface{} {
	if h.data.Size() == 0 {
		panic("Can not findMax when heap is empty.")
	}
	return h.data.Get(0)
}

// 取出堆中最大元素
func (h *maxHeap) ExtractMax() interface{} {
	ret := h.FindMax()

	h.data.Swap(0, h.data.Size()-1)
	h.data.RemoveLast()
	h.siftDown(0)

	return ret
}

func (h *maxHeap) ExtractMin() interface{} {
	ret := h.FindMin()

	h.data.RemoveLast()

	return ret
}

// data[j] 是 leftChild 和 rightChild 中的最大值
func (h *maxHeap) siftDown(k int) {
	for leftChild(k) < h.data.Size() {
		j := leftChild(k)
		// j+1是右孩子索引，如果存在右孩子比较后获得左右孩子中较大值的索引
		ret, err := common.Compare(h.data.Get(j+1), h.data.Get(j))
		if j+1 < h.data.Size() && err == nil && ret > 0 {
			j++
		}
		// data[j] 是 leftChild 和 rightChild 中的最大值
		if ret, err := common.Compare(h.data.Get(k), h.data.Get(j)); err == nil && ret > 0 {
			break
		}

		h.data.Swap(k, j)
		k = j
	}
}

// 取出堆中的最大元素，并且替换成元素e
func (h *maxHeap) Replace(e interface{}) interface{} {
	ret := h.FindMax()

	h.data.Set(0, e)
	h.siftDown(0)

	return ret
}
