package list

import (
	"bytes"
	"fmt"
)

// arrayStack 栈

type IStack interface {
	Size() int
	IsEmpty() bool
	Push(interface{})
	Pop() interface{}
	Peek() interface{}
}

type arrayStack struct {
	array *ArrayList
}

func NewArrayStack(capacity int) *arrayStack {
	return &arrayStack{
		array: NewArrayList(capacity),
	}
}

func (s *arrayStack) Size() int {
	return s.array.Size()
}

func (s *arrayStack) IsEmpty() bool {
	return s.array.IsEmpty()
}

// 压入栈顶元素
func (s *arrayStack) Push(element interface{}) {
	s.array.AddLast(element)
}

// 弹出栈顶元素
func (s *arrayStack) Pop() interface{} {
	return s.array.RemoveLast()
}

// 查看栈顶元素
func (s *arrayStack) Peek() interface{} {
	return s.array.Get(s.Size() - 1)
}

func (s *arrayStack) String() string {
	var buffer bytes.Buffer

	buffer.WriteString("arrayStack: ")
	buffer.WriteString("[")
	for i := 0; i < s.array.Size(); i++ {
		buffer.WriteString(fmt.Sprint(s.array.Get(i)))
		if i != s.array.Size()-1 {
			buffer.WriteString(", ")
		}
	}
	buffer.WriteString("] top")

	return buffer.String()
}

// 判断某些符号是否对称
func IsValid(s string, brackets map[rune]rune) bool {
	stack := NewArrayStack(20)

	for _, char := range s {
		if char == '(' || char == '{' || char == '[' {
			// 入栈
			stack.Push(char)
		} else if stack.Size() > 0 && brackets[char] == stack.Peek() {
			// 栈中有数据，且此元素与栈尾元素相同
			stack.Pop()
		} else {
			return false
		}
	}

	// 循环结束，栈中还有数据则 false
	return stack.Size() == 0
}
