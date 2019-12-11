package queue

import (
	"fmt"
	"testing"
)

func TestArrayQueue(t *testing.T) {
	q := NewArrayQueue(10)
	q.Enqueue("1")
	q.Enqueue("2")
	fmt.Println(q.Dequeue())
}
