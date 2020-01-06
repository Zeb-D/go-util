package todo

import (
	"fmt"
	"testing"
	"unsafe"
)

func TestSize(t *testing.T) {
	// https://wizardforcel.gitbooks.io/gopl-zh/ch13/ch13-01.html
	// Go語言的規范併沒有要求一個字段的聲明順序和內存中的順序是一致的，所以理論上一個編譯器可以隨意地重新排列每個字段的內存位置，隨然在寫作本書的時候編譯器還沒有這麽做。
	// 下面的結構體雖然有着相同的字段，但是第一種寫法比另外的需要多50%的內存。
	t1 := T1{}
	t2 := T2{}
	fmt.Println(unsafe.Sizeof(t1))
	fmt.Println(unsafe.Sizeof(t2))

	// unsafe.Alignof 函數返迴對應參數的類型需要對齊的倍數. 和 Sizeof 類似, Alignof 也是返迴一個常量表達式,
	// 對應一個常量.
	// 通常情況下布爾和數字類型需要對齊到它們本身的大小(最多8個字節), 其它的類型對齊到機器字大小.
	fmt.Println(unsafe.Alignof(t1))
	fmt.Println(unsafe.Alignof(t2))

	// unsafe.Offsetof 函數的參數必鬚是一個字段 x.f, 然後返迴 f 字段相對於 x 起始地址的偏移量, 包括可能的空洞.
	fmt.Println(unsafe.Offsetof(t1.f))
	fmt.Println(unsafe.Offsetof(t1.b))
}

type T1 struct {
	b bool
	i int64
	f float32
}

type T2 struct {
	b bool
	f float32
	i int64
}

type Person struct {
	name string
	age  int
}
type APerson struct {
	Person
	sex int
	int
}

func TestPerson(t *testing.T) {
	var ap = APerson{
		sex: 1,
	}
	ap.name = "p"
	ap.int = 11
	println(ap.name, "->", ap.age, "->", ap.int)
	// ap instance Person is false
}

func Type(a interface{}) {
	t, ok := a.(Person)
	println(t.age, "->", ok)
}
