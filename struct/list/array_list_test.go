package list

import (
	"fmt"
	"testing"
)

func TestArrayList_Contains(t *testing.T) {
	array := NewArrayList(10)
	array.AddFirst(12121)
	array.AddLast(12122)
	array.AddFirst(12121)
	fmt.Println(array)
	fmt.Println(array.Contains(12122))
	fmt.Println(array.Contains("12122"))
	fmt.Println(array.Remove(1))
	fmt.Println(array.Find(12122))
	fmt.Println(array.RemoveElement(123))
	fmt.Println(array.RemoveAllElement(12121))
}
