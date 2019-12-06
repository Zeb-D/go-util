package list

import "fmt"

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
	//Node := &Node{e: e}
	//Node.next = linkedList.head
	//linkedList.head = Node
	linkedList.head = &Node{e, linkedList.head}
	linkedList.size++
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
// 在链表中不是一个常用的操作，练习用：）
func (linkedList *LinkedList) Add(index int, e interface{}) {
	if index < 0 || index > linkedList.size {
		panic("Add failed. Illegal index.")
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
}

// 在链表末尾添加新的元素e
func (linkedList *LinkedList) AddLast(e interface{}) {
	linkedList.Add(linkedList.size, e)
}
