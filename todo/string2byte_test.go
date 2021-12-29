package todo

import (
	"fmt"
	"strings"
	"testing"
	"unsafe"
)

func String2bytes(s string) []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&s))
	h := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&h))
}

func Bytes2String(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

var ss = strings.Repeat("hello", 1024)

func testDefault() {
	a := []byte(ss)
	_ = string(a)
}

func testUnsafe() {
	a := String2bytes(ss)
	_ = Bytes2String(a)
}

//	BenchmarkTestDefault-4   	  670196	      1736 ns/op
func BenchmarkTestDefault(b *testing.B) {
	for i := 0; i < b.N; i++ {
		testDefault()
	}
}

//	BenchmarkTestUnsafe-4   	1000000000	         0.570 ns/op
func BenchmarkTestUnsafe(b *testing.B) {
	for i := 0; i < b.N; i++ {
		testUnsafe()
	}
}

func TestString2bytes(t *testing.T) {
	bs := String2bytes(ss)
	fmt.Println(bs)
	fmt.Println(Bytes2String(bs))

	println("string Pointer")
	x := (*[2]uintptr)(unsafe.Pointer(&ss))
	fmt.Println(unsafe.Pointer(&ss))
	fmt.Println(x)
	fmt.Println(len(ss))

	bb := *(*[]byte)(unsafe.Pointer(&[2]uintptr{x[0], 5})) // 起始地址，加长度
	fmt.Println(string(bb))
	fmt.Println(bb)
	fmt.Println(*(*string)(unsafe.Pointer(&bb)))
	fmt.Println(*(*byte)(unsafe.Pointer(&bb)))

	println("int Pointer")
	length := len(ss)
	fmt.Println(unsafe.Pointer(&length))
	fmt.Println((*[2]uintptr)(unsafe.Pointer(&length)))

	xx := (*[2]uintptr)(unsafe.Pointer(&length))
	bbb := *(*int)(unsafe.Pointer(&[2]uintptr{xx[0], 2}))
	fmt.Println(bbb)
}
