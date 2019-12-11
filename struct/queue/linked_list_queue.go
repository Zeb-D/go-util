package queue

import (
	"bytes"
	"fmt"
)

type node struct {
	e    interface{}
	next *node
}

func (n *node) String() string {
	return fmt.Sprint(n.e)
}

type LinkedListQueue struct {
	head *node
	tail *node
	size int
}

func NewLinkedListQueue() *LinkedListQueue {
	return &LinkedListQueue{}
}

func (linkedListQueue *LinkedListQueue) Size() int {
	return linkedListQueue.size
}

func (linkedListQueue *LinkedListQueue) IsEmpty() bool {
	return linkedListQueue.size == 0
}

func (linkedListQueue *LinkedListQueue) Enqueue(e interface{}) {
	if linkedListQueue.tail == nil {
		linkedListQueue.tail = &node{e: e}
		linkedListQueue.head = linkedListQueue.tail
	} else {
		linkedListQueue.tail.next = &node{e: e}
		linkedListQueue.tail = linkedListQueue.tail.next
	}
	linkedListQueue.size++
}

func (linkedListQueue *LinkedListQueue) Dequeue() interface{} {
	if linkedListQueue.IsEmpty() {
		panic("Cannot dequeue from an empty queue.")
	}
	head := linkedListQueue.head
	linkedListQueue.head = linkedListQueue.head.next
	head.next = nil
	if head == nil {
		linkedListQueue.tail = nil
	}
	linkedListQueue.size--

	return head.e
}

func (linkedListQueue *LinkedListQueue) Front() interface{} {
	if linkedListQueue.IsEmpty() {
		panic("Cannot dequeue from an empty queue.")
	}

	return linkedListQueue.head.e
}

func (linkedListQueue *LinkedListQueue) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("Queue: front ")

	cur := linkedListQueue.head
	for cur != nil {
		buffer.WriteString(fmt.Sprintf("%v ->", cur.e))
		cur = cur.next
	}
	buffer.WriteString("NULL tail")
	return buffer.String()
}
