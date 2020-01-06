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

//example Size: 8
//Alignment Boundary: 8
//BoolValue = Size: 1 Offset: 0 Addr: 0xc0001818a0
//IntValue = Size: 2 Offset: 2 Addr: 0xc0001818a2
//FloatValue = Size: 4 Offset: 4 Addr: 0xc0001818a4
//Next = Size: 1 Offset: 0 Addr: 0xc0001818a8
func TestUnsafe(t *testing.T) {
	example := &Example{
		BoolValue:  true,
		IntValue:   10,
		FloatValue: 3.141592,
	}

	exampleNext := &Example{
		BoolValue:  true,
		IntValue:   10,
		FloatValue: 3.141592,
	}

	alignmentBoundary := unsafe.Alignof(example)

	sizeBool := unsafe.Sizeof(example.BoolValue)
	offsetBool := unsafe.Offsetof(example.BoolValue)

	sizeInt := unsafe.Sizeof(example.IntValue)
	offsetInt := unsafe.Offsetof(example.IntValue)

	sizeFloat := unsafe.Sizeof(example.FloatValue)
	offsetFloat := unsafe.Offsetof(example.FloatValue)

	sizeBoolNext := unsafe.Sizeof(exampleNext.BoolValue)
	offsetBoolNext := unsafe.Offsetof(exampleNext.BoolValue)

	fmt.Printf("example Size: %d\n", unsafe.Sizeof(example))

	fmt.Printf("Alignment Boundary: %d\n", alignmentBoundary)

	fmt.Printf("BoolValue = Size: %d Offset: %d Addr: %v\n",
		sizeBool, offsetBool, &example.BoolValue)

	fmt.Printf("IntValue = Size: %d Offset: %d Addr: %v\n",
		sizeInt, offsetInt, &example.IntValue)

	fmt.Printf("FloatValue = Size: %d Offset: %d Addr: %v\n",
		sizeFloat, offsetFloat, &example.FloatValue)

	fmt.Printf("Next = Size: %d Offset: %d Addr: %v\n",
		sizeBoolNext, offsetBoolNext, &exampleNext.BoolValue)

}
