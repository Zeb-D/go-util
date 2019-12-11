package list

import (
	"fmt"
	"testing"
)

func TestLinkedList_Add(t *testing.T) {
	list := &LinkedList{}
	fmt.Println(list.IsEmpty())
	list.AddFirst("aaa")
	fmt.Println(list.Size())
	list.AddFirst(111)
	fmt.Println(list.head)
	list.RemoveFirst()
	fmt.Println(list.head)
	fmt.Println(list.size)
	list.AddFirst(1112)
	list.AddFirst(1113)
	list.Add(1, 222)
	fmt.Println(list)
	fmt.Println(list.GetFirst())
	fmt.Println(list.Get(1))
	fmt.Println(list.RemoveFirst())
	fmt.Println(list.GetLast())
}

func TestLinkedList_Remove(t *testing.T) {
	list := &LinkedList{}
	list.AddFirst("pppp")
	list.AddFirst("pppp1")
	list.Add(1, "aaa")
	fmt.Println(list)
	fmt.Println(list.Remove(2))
	fmt.Println(list.Remove(1))
	fmt.Println(list.RemoveFirst())
	fmt.Println(list)
}
