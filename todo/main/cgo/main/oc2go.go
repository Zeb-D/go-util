package main

//
//import "C"
//
////export helloV5
//func helloV5(value string) *C.char { // 如果函数有返回值，则要将返回值转换为C语言对应的类型
//	return C.CString("helloV5" + value)
//}
//
////C调用到Go这种情况比较复杂，Go一般是便以为c-shared/c-archive的库给C调用。
//func main() {
//	// 此处一定要有main函数，有main函数才能让cgo编译器去把包编译成C的库
//}
////如果Go函数有多个返回值，会生成一个C结构体进行返回，结构体定义参考生成的.h文件
////go build -buildmode=c-shared -o oc2go.so oc2go.go
