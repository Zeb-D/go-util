package queue

import (
	"bytes"
	"fmt"
)

// LoopQueue 循环队列
type LoopQueue struct {
	data  []interface{}
	front int
	tail  int
	size  int
}

func NewLoopQueue(capacity int) *LoopQueue {
	// 由于不浪费空间，所以data静态数组的大小是capacity
	// 而不是capacity + 1
	return &LoopQueue{
		data: make([]interface{}, capacity),
	}
}

func (loopQueue *LoopQueue) Capacity() int {
	return len(loopQueue.data)
}

func (loopQueue *LoopQueue) Size() int {
	return loopQueue.size
}

func (loopQueue *LoopQueue) IsEmpty() bool {
	// 注意，我们不再使用front和tail之间的关系来判断队列是否为空，而直接使用size
	return loopQueue.size == 0
}

// 入队
func (loopQueue *LoopQueue) Enqueue(e interface{}) {
	// 注意，我们不再使用front和tail之间的关系来判断队列是否为满，而直接使用size
	capacity := loopQueue.Capacity()
	if loopQueue.size == capacity {
		loopQueue.resize(capacity * 2)
	}
	loopQueue.data[loopQueue.tail] = e
	loopQueue.tail = (loopQueue.tail + 1) % len(loopQueue.data)
	loopQueue.size++
}

// 获得队列头部元素
func (loopQueue *LoopQueue) Dequeue() (e interface{}) {
	if loopQueue.IsEmpty() {
		panic("Cannot dequeue from empty queue")
	}

	e = loopQueue.data[loopQueue.front]
	loopQueue.data[loopQueue.front] = nil
	// 循环队列需要执行求余运算
	loopQueue.front = (loopQueue.front + 1) % len(loopQueue.data)
	loopQueue.size--
	if loopQueue.size == loopQueue.Capacity()/4 && loopQueue.size != 0 {
		loopQueue.resize(loopQueue.Capacity() / 2)
	}

	return
}

// 查看队列头部元素
func (loopQueue *LoopQueue) Front() interface{} {
	if loopQueue.IsEmpty() {
		panic("Queue is empty")
	}

	return loopQueue.data[loopQueue.front]
}

func (loopQueue *LoopQueue) resize(capacity int) {
	length := len(loopQueue.data)
	newData := make([]interface{}, capacity)
	for i := 0; i < loopQueue.size; i++ {
		newData[i] = loopQueue.data[(i+loopQueue.front)%length]
	}
	loopQueue.data = newData
	loopQueue.front = 0
	loopQueue.tail = loopQueue.size
}

func (loopQueue *LoopQueue) String() string {
	var buffer bytes.Buffer

	buffer.WriteString(fmt.Sprintf("Queue: size = %d, capacity = %d\n", loopQueue.size, loopQueue.Capacity()))
	buffer.WriteString("front [")
	for i := loopQueue.front; i != loopQueue.tail; i = (i + 1) % len(loopQueue.data) {
		// fmt.Sprint 将 interface{} 类型转换为字符串
		buffer.WriteString(fmt.Sprint(loopQueue.data[i]))
		if (i+1)%len(loopQueue.data) != loopQueue.tail {
			buffer.WriteString(", ")
		}
	}
	buffer.WriteString("] tail")

	return buffer.String()
}
