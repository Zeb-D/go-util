package others

import (
	"fmt"
	"reflect"
	"testing"
	"unsafe"
)

//主要研究string 如何零拷贝出一个[]byte，注意这个切片不能修改

func TestString2Byte(t *testing.T) {
	a := "aaa"
	bs := []byte(a)
	fmt.Printf("%v \n", bs)
}

func TestStringNoCopy2Byte(t *testing.T) {
	a := "aaa"
	ssh := *(*reflect.StringHeader)(unsafe.Pointer(&a))
	b := *(*[]byte)(unsafe.Pointer(&ssh))
	fmt.Printf("%v \n", b)
}

func TestStringClone(t *testing.T) {
	s := "abcdefghijklmn"
	s1 := s[:4]

	sHeader := (*reflect.StringHeader)(unsafe.Pointer(&s))
	s1Header := (*reflect.StringHeader)(unsafe.Pointer(&s1))
	fmt.Println(sHeader.Len == s1Header.Len)
	fmt.Println(sHeader.Data == s1Header.Data)

	// Output:
	// false
	// true
	fmt.Println(Clone(s1))
}

func Clone(s string) string {
	if len(s) == 0 {
		return ""
	}
	b := make([]byte, len(s))
	copy(b, s)
	return *(*string)(unsafe.Pointer(&b))
}

func String2Bytes(s string) []byte {
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	bh := reflect.SliceHeader{
		Data: sh.Data,
		Len:  sh.Len,
		Cap:  sh.Len,
	}
	return *(*[]byte)(unsafe.Pointer(&bh))
}

//https://mp.weixin.qq.com/s/lHk2L8p8ClmXXuf4zS0qLQ
//GOTRACEBACK=system go run main.go
//如果想获取 core dump 文件，那么就应该把 GOTRACEBACK 的值设置为 crash 。
// 当然，我们还可以通过 runtime/debug 包中的 SetTraceback 方法来设置堆栈打印级别。
//mac 系统下的 Go 限制了生成 core dump 文件，这个在 Go 源码 src/runtime/signal_unix.go 中有相关说明。
func Modify() {
	a := "hello"
	b := String2Bytes(a)
	b[0] = 'H' //b是引用了a的byte指针，底层不允许变更，会报错
}
