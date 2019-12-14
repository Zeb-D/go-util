package Map

import "fmt"

type LinkedListNode struct {
	key  interface{}
	val  interface{}
	next *LinkedListNode
}

type LinkedListMap struct {
	head *LinkedListNode
	size int
}

func NewLinkedListMap() *LinkedListMap {
	return &LinkedListMap{
		head: &LinkedListNode{},
	}
}

func (l *LinkedListMap) Size() int {
	return l.size
}

func (l *LinkedListMap) IsEmpty() bool {
	return l.size == 0
}

func (l *LinkedListMap) getNode(key interface{}) *LinkedListNode {
	prev := l.head.next
	for prev != nil {
		if prev.key == key {
			return prev
		}
		prev = prev.next
	}

	return nil
}

func (l *LinkedListMap) Add(key interface{}, val interface{}) {
	n := l.getNode(key)

	if n == nil {
		newNode := &LinkedListNode{
			key:  key,
			val:  val,
			next: l.head.next,
		}
		l.head.next = newNode
		l.size++
	} else {
		n.val = val
	}
}

func (l *LinkedListMap) Get(key interface{}) interface{} {
	n := l.getNode(key)
	if n == nil {
		return nil
	} else {
		return n.val
	}
}

func (l *LinkedListMap) Set(key interface{}, val interface{}) {
	Node := l.getNode(key)
	if Node == nil {
		panic(fmt.Sprintf("%v, doesn't exist", key))
	}

	Node.val = val
}

func (l *LinkedListMap) Remove(key interface{}) interface{} {
	prev := l.head
	for prev.next != nil {
		if prev.next.key == key {
			delNode := prev.next
			prev.next = delNode.next
			delNode.next = nil

			return delNode.val
		}
		prev = prev.next
	}

	return nil
}

func (l *LinkedListMap) Contains(key interface{}) bool {
	return l.getNode(key) != nil
}
