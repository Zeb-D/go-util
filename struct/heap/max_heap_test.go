package heap

import (
	"fmt"
	"testing"
)

func TestMaxHeap(t *testing.T) {
	heap := NewMaxHeap(20)
	heap.Add(22)
	heap.Add(11)
	fmt.Println("1->", heap)
	heap.Add(21)
	heap.Add(11)
	fmt.Println("2->", heap)
}
