package queue

import (
	"fmt"
	"testing"
)

func TestLinkedListQueue(t *testing.T) {
	q := NewLinkedListQueue()
	q.Enqueue(1)
	fmt.Println(q.size)
	q.Enqueue("121")
	fmt.Println(q.Front())
	fmt.Println(q.IsEmpty())
	fmt.Println(q.Dequeue())
	fmt.Println(q)
}
