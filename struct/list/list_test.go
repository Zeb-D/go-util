package list

import (
	"fmt"
	"testing"
)

func TestLoopLinkedList(t *testing.T) {
	l := NewLoopLinkedList()
	l.AddFirst("a")
	l.Add(1, 2)
	fmt.Println(l)
	l.RemoveFirst()
	fmt.Println(l)
}

func TestTwoWayLinkedList(t *testing.T) {
	l := NewTwoWayLinkedList()
	l.AddFirst("a")
	l.Add(1, 2)
	fmt.Println(l)
	l.RemoveFirst()
	fmt.Println(l)
}
