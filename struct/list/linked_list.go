package list

import (
	"errors"
	"fmt"
)

type Node struct {
	e    interface{}
	next *Node
}

func (n *Node) String() string {
	return fmt.Sprintf("[e:%v,next:%v]", n.e, n.next)
}

func (n *Node) Next() *Node {
	return n.next
}

func (n *Node) Value() interface{} {
	return n.e
}

type LinkedList struct {
	head *Node
	size int
}

func (linkedList *LinkedList) String() string {
	return fmt.Sprintf("[size:%v,head:%v]", linkedList.size, linkedList.head)
}

func (linkedList *LinkedList) Head() *Node {
	return linkedList.head
}

// 获取链表中的元素个数
func (linkedList *LinkedList) GetSize() int {
	return linkedList.size
}

// 返回链表是否为空
func (linkedList *LinkedList) IsEmpty() bool {
	return linkedList.size == 0
}

// 在链表头添加新的元素e
func (linkedList *LinkedList) AddFirst(e interface{}) {
	linkedList.head = &Node{e, linkedList.head}
	linkedList.size++
}

// 从链表中删除index(0-based)位置的元素，返回删除的元素
func (linkedList *LinkedList) Remove(index int) (interface{}, error) {
	if index < 0 || index >= linkedList.size {
		return nil, errors.New("remove failed. Index is illegal.")
	}

	if linkedList == nil || linkedList.head == nil {
		return nil, nil
	}
	// 第一个节点
	if index == 0 {
		e := linkedList.head.e

		curlNode := linkedList.head
		linkedList.head = curlNode.next
		curlNode.next = nil
		linkedList.size--

		return e, nil
	}

	// prev 是待删除元素的前一个元素
	prevNode := linkedList.head
	for i := 0; i < index-1; i++ {
		prevNode = prevNode.next
	}

	curlNode := prevNode.next
	e := curlNode.e
	prevNode.next = curlNode.next
	curlNode.next = nil

	linkedList.size--

	return e, nil
}

func (linkedList *LinkedList) RemoveFirst() interface{} {
	if linkedList == nil || linkedList.head == nil {
		return nil
	}
	e := linkedList.head.e
	linkedList.head = linkedList.head.next
	linkedList.size--
	return e
}

// 在链表的index(0-based)位置添加新的元素e
func (linkedList *LinkedList) Add(index int, e interface{}) error {
	if index < 0 || index > linkedList.size {
		return errors.New("add failed. illegal index.")
	}

	if index == 0 {
		linkedList.AddFirst(e)
	} else {
		// 获得待插入节点的前一个节点
		prev := linkedList.head
		for i := 0; i < index-1; i++ {
			prev = prev.next
		}

		// 插入新节点
		//Node := &Node{e: e, next: prev.next}
		//prev.next = Node
		prev.next = &Node{e, prev.next}
		linkedList.size++
	}
	return nil
}

// 在链表末尾添加新的元素e
func (linkedList *LinkedList) AddLast(e interface{}) error {
	return linkedList.Add(linkedList.size, e)
}

// 修改链表的第index(0-based)个位置的元素为e
func (linkedList *LinkedList) Set(index int, e interface{}) error {
	if index < 0 || index >= linkedList.size {
		return errors.New("Set failed. Illegal index.")
	}

	cur := linkedList.head
	for i := 0; i < index; i++ {
		cur = cur.next
	}
	cur.e = e
	return nil
}

// 查找链表是否存在元素e
func (linkedList *LinkedList) Contains(e interface{}) bool {
	cur := linkedList.head

	for cur != nil {
		if cur.e == e {
			return true
		}
		cur = cur.next
	}
	return false
}

// 获得链表的第index(0-based)个位置的元素
func (linkedList *LinkedList) Get(index int) (interface{}, error) {
	if index < 0 || index >= linkedList.size {
		return nil, errors.New("add failed. illegal index.")
	}
	count := 0
	for node := linkedList.head; node != nil; node = node.next {
		if count == index {
			return node.e, nil
		}
		count++
	}
	return nil, nil
}

// 获得链表的第一个元素
func (linkedList *LinkedList) GetFirst() (interface{}, error) {
	return linkedList.Get(0)
}

// 获得链表的最后一个元素
func (linkedList *LinkedList) GetLast() (interface{}, error) {
	return linkedList.Get(linkedList.size - 1)
}
