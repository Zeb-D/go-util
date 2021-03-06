package todo

import (
	"fmt"
	"testing"
	"unsafe"
)

type Example struct {
	BoolValue  bool
	IntValue   int16
	FloatValue float32
}

var example = &Example{
	BoolValue:  true,
	IntValue:   10,
	FloatValue: 3.141592,
}

var exampleNext = &Example{
	BoolValue:  true,
	IntValue:   10,
	FloatValue: 3.141592,
}

//example Size: 8
//Alignment Boundary: 8
//BoolValue = Size: 1 Offset: 0 Addr: 0xc0001818a0
//IntValue = Size: 2 Offset: 2 Addr: 0xc0001818a2
//FloatValue = Size: 4 Offset: 4 Addr: 0xc0001818a4
//Next = Size: 1 Offset: 0 Addr: 0xc0001818a8
//在BoolValue和IntValue字段之间填充1个字节。
// 偏移值和两个地址之间的差异是2个字节。
// 您还可以看到下一个内存分配是从结构中的最后一个字段开始4个字节。
func TestUnsafe_alignof(t *testing.T) {
	//alignment（数据结构对齐）简单来说，它让数据在内存中以某种的布局来存放，使该数据的读取能够更加的快速。
	// 这个接收一个变量作为参数，并返回这个变量的对齐字节。
	alignmentBoundary := unsafe.Alignof(example)

	//由于 alignment 机制的要求，n2 的内存起始地址应该是自身大小的整数倍，
	// 也就是说它的起始地址只能是 0、2、4、6、8 等偶数，所以 n2 的起始地址没有紧接着 n1 后面，而是空出了 1 个字节。
	fmt.Printf("Alignment Boundary: %d\n", alignmentBoundary)

	//因为 alignment 除了要求字段的其实地址应该是自身大小的整数倍，
	// 还要求整个结构体的大小，是结构体中最大的字段的大小的整数倍，这使得结构体可以由多个内存块组成，
	// 其中每个内存块的大小都等于最大的字段的大小。
}

func TestSizeof(t *testing.T) {
	fmt.Printf("example Size: %d\n", unsafe.Sizeof(example))

	sizeBool := unsafe.Sizeof(example.BoolValue)
	offsetBool := unsafe.Offsetof(example.BoolValue)

	sizeInt := unsafe.Sizeof(example.IntValue)
	offsetInt := unsafe.Offsetof(example.IntValue)

	sizeFloat := unsafe.Sizeof(example.FloatValue)
	offsetFloat := unsafe.Offsetof(example.FloatValue)

	fmt.Printf("BoolValue = Size: %d Offset: %d Addr: %v\n",
		sizeBool, offsetBool, &example.BoolValue)

	fmt.Printf("IntValue = Size: %d Offset: %d Addr: %v\n",
		sizeInt, offsetInt, &example.IntValue)

	fmt.Printf("FloatValue = Size: %d Offset: %d Addr: %v\n",
		sizeFloat, offsetFloat, &example.FloatValue)

	sizeBoolNext := unsafe.Sizeof(exampleNext.BoolValue)
	offsetBoolNext := unsafe.Offsetof(exampleNext.BoolValue)

	fmt.Printf("Next = Size: %d Offset: %d Addr: %v\n",
		sizeBoolNext, offsetBoolNext, &exampleNext.BoolValue)
}

func TestMap(t *testing.T) {
	var m map[string]int
	m = map[string]int{"one": 1, "two": 2}
	n := m
	fmt.Printf("%p\n", &m) //0xc000074018
	fmt.Printf("%p\n", &n) //0xc000074020
	fmt.Println(m)         // map[two:2 one:1]
	fmt.Println(n)         //map[one:1 two:2]
	changeMap(&m)
	fmt.Printf("%p\n", &m) //0xc000074018
	fmt.Printf("%p\n", &n) //0xc000074020
	fmt.Println(m)         //map[one:1 two:2 three:3]
	fmt.Println(n)         //map[two:2 three:3 one:1]

	var a = 75.0
	var p1 = &a
	var p2 = &a

	if p1 == p2 {
		fmt.Println("Both pointers p1 and p2 point to the same variable.")
	}

	a = 7.98
	var p = &a
	var pp = &p

	fmt.Println("a = ", a)
	fmt.Println("address of a = ", &a)

	fmt.Println("p = ", p)
	fmt.Println("address of p = ", &p)

	fmt.Println("pp = ", pp)

	// Dereferencing a pointer to pointer
	fmt.Println("*pp = ", *pp)
	fmt.Println("**pp = ", **pp)

	a = 1000
	p = &a

	fmt.Println("a (before) = ", a)

	// Changing the value stored in the pointed variable through the pointer
	*p = 2000

	fmt.Println("a (after) = ", a)

	var p3 *int

	p3 = new(int)
	*p3 = 1
	fmt.Println("p3=", *p3)
}

func changeMap(m *map[string]int) {
	//m["three"] = 3 //这种方式会报错 invalid operation: m["three"] (type *map[string]int does not support indexing)
	(*m)["three"] = 3                    //正确
	fmt.Printf("changeMap func %p\n", m) //changeMap func 0x0
}

//unsafe 包压根就没有源码[5]呀。它有函数的签名和类型定义，但是没有实现的代码：无论是 Go 还是汇编的代码都没有。
// 之所以会出现这个情况，是因为 unsafe 包的功能需要在层次更低的编译器层面实现，所以这个包其实是内置在编译器里面实现的，
// 这个 .go 文件只是为了达到文档记录的目的。
// 所以我在上文反复强调要严格遵循 Go 研发团队的文档和建议，因为你也只能看到这些文档。
