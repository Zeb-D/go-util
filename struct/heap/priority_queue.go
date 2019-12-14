package heap

type PriorityQueue struct {
	maxHeap *maxHeap
}

func NewPriorityQueue(capacity int) *PriorityQueue {
	return &PriorityQueue{NewMaxHeap(capacity)}
}

func (q *PriorityQueue) Size() int {
	return q.maxHeap.Size()
}

func (q *PriorityQueue) IsEmpty() bool {
	return q.maxHeap.IsEmpty()
}

func (q *PriorityQueue) Enqueue(e interface{}) {
	q.maxHeap.Add(e)
}

func (q *PriorityQueue) Dequeue() interface{} {
	return q.maxHeap.ExtractMax()
}

func (q *PriorityQueue) Front() interface{} {
	return q.maxHeap.FindMax()
}
