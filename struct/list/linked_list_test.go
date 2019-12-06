package list

import (
	"fmt"
	"testing"
)

func TestLinkedList_Add(t *testing.T) {
	list := LinkedList{}
	fmt.Println(list.IsEmpty())
	list.AddFirst("aaa")
	fmt.Println(list.GetSize())
	list.AddFirst(111)
	fmt.Println(list.head)
	list.RemoveFirst()
	fmt.Println(list.head)
	fmt.Println(list.size)
}
