//go2cpp.go
package main

/*
#include "go2cpp.h"
int SayHelloV3();
*/
import "C"
import (
	"fmt"
)

//编译失败，要使用// go run ../main
//CGO提供的这种面向C语言接口的编程方式，
//使得开发者可以使用是任何编程语言来对接口进行实现，只要最终满足C语言接口即可。
func main() {
	ret := C.SayHelloV3()
	fmt.Println(ret)
}
