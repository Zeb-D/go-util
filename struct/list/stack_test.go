package list

import (
	"fmt"
	"testing"
)

func TestStack_Pop(t *testing.T) {
	stack := NewArrayStack(10)

	fmt.Println(stack)
	for i := 0; i < 10; i++ {
		stack.Push(i)
		fmt.Println(stack)
	}

	fmt.Println(stack.Pop())
	fmt.Println(stack)
}

func TestNewStack(t *testing.T) {
	brackets := map[rune]rune{')': '(', ']': '[', '}': '{'}
	fmt.Println(IsValid("()[]{}", brackets))
	fmt.Println(IsValid("([)]", brackets))
}
