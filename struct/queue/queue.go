package queue

// Queue  先进先出
type IQueue interface {
	Size() int
	IsEmpty() bool
	Enqueue(interface{})
	Dequeue() interface{}
	Front() interface{}
}
