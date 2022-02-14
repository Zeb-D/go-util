package main

//#include <stdio.h>        //  序文中可以链接标准C程序库
import "C" // import "C"更像是一个关键字，CGO工具在预处理时会删掉这一行

//通过import“C”语句启用CGO特性后，CGO会将上一行代码所处注释块的内容视为C代码块，
//被称为序文（preamble）。
//https://mp.weixin.qq.com/s/I9IvJnKeJY8IpVB4lcWXNA
//func main() {
//	C.puts(C.CString("Hello, Cgo\n"))
//}

//go build -x cgo.go
