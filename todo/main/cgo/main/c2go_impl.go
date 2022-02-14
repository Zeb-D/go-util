//c2go_impl.go
package main

//#include <c2go.h>
import "C"
import "fmt"

//export SayHelloV4
func SayHelloV4(str *C.char) {
	fmt.Println(C.GoString(str))
}

//CGO的//export SayHello指令将Go语言实现的SayHello函数导出为C语言函数。
