package set

import "github.com/Zeb-D/go-util/struct/list"

// base on linked list
type LinkedListSet struct {
	LinkedList *list.LinkedList
}

func NewLinkedListSet() *LinkedListSet {
	return &LinkedListSet{
		LinkedList: list.NewLinkedList(),
	}
}

func (l *LinkedListSet) Add(e interface{}) {
	if !l.LinkedList.Contains(e) {
		l.LinkedList.AddFirst(e)
	}
}

func (l *LinkedListSet) Remove(e interface{}) {
	l.LinkedList.RemoveElement(e)
}

func (l *LinkedListSet) Contains(e interface{}) bool {
	return l.LinkedList.Contains(e)
}

func (l *LinkedListSet) Size() int {
	return l.LinkedList.Size()
}

func (l *LinkedListSet) IsEmpty() bool {
	return l.LinkedList.IsEmpty()
}
