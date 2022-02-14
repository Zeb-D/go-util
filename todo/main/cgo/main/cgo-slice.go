//cgo-slice.go
package main

//因此Go的切片不能直接传递给C使用，而是需要取切片的内部缓冲区的首地址(即首个元素的地址)来传递给C使用。
// 使用这种方式把Go的内存空间暴露给C使用，可以大大减少Go和C之间参数传递时内存拷贝的消耗。

/*
int SayHelloV6(char* buff, int len) {
    char hello[] = "Hello Cgo!";
    int movnum = len < sizeof(hello) ? len:sizeof(hello);
    memcpy(buff, hello, movnum);                        // go字符串没有'\0'，所以直接内存拷贝
    return movnum;
}

*/
import "C"
import (
	"fmt"
	"unsafe"
)

func main() {
	buff := make([]byte, 8)
	C.SayHelloV6((*C.char)(unsafe.Pointer(&buff[0])), C.int(len(buff)))
	a := string(buff)
	fmt.Println(a)
}
