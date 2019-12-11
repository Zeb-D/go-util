package queue

import (
	"bytes"
	"fmt"
	"github.com/Zeb-D/go-util/struct/list"
)

// 数组队列的局限性：出队列的操作时间复杂度为 n
type ArrayQueue struct {
	array *list.ArrayList
}

func NewArrayQueue(capacity int) *ArrayQueue {
	return &ArrayQueue{
		array: list.NewArrayList(capacity),
	}
}

func (arrayQueue *ArrayQueue) Size() int {
	return arrayQueue.array.Size()
}

func (arrayQueue *ArrayQueue) IsEmpty() bool {
	return arrayQueue.array.IsEmpty()
}

func (arrayQueue *ArrayQueue) Capacity() int {
	return arrayQueue.array.Capacity()
}

func (arrayQueue *ArrayQueue) Enqueue(e interface{}) {
	arrayQueue.array.AddLast(e)
}

func (arrayQueue *ArrayQueue) Dequeue() interface{} {
	return arrayQueue.array.RemoveFirst()
}

func (arrayQueue *ArrayQueue) Front() interface{} {
	return arrayQueue.array.Get(0)
}

func (arrayQueue *ArrayQueue) String() string {
	var buffer bytes.Buffer

	buffer.WriteString("Queue: ")
	buffer.WriteString("front [")
	for i := 0; i < arrayQueue.array.Size(); i++ {
		// fmt.Sprint 将 interface{} 类型转换为字符串
		buffer.WriteString(fmt.Sprint(arrayQueue.array.Get(i)))
		if i != (arrayQueue.array.Size() - 1) {
			buffer.WriteString(", ")
		}
	}
	buffer.WriteString("] tail")

	return buffer.String()
}
