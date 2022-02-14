package main

/*
#cgo LDFLAGS: -L/usr/local/lib

#include <stdio.h>
#include <stdlib.h>
#define REPEAT_LIMIT 2              // CGO会保留C代码块中的宏定义
typedef struct{                     // 自定义结构体
    int repeat_time;
    char* str;
}blob;
int SayHello0(blob* pblob) {  // 自定义函数
    for ( ;pblob->repeat_time < REPEAT_LIMIT; pblob->repeat_time++){
        puts(pblob->str);
    }
    return 0;
}
*/
import "C"

//import (
//	"fmt"
//	"unsafe"
//)
//
//func main() {
//	cblob := C.blob{} // 在GO程序中创建的C对象，存储在Go的内存空间
//	cblob.repeat_time = 0
//
//	cblob.str = C.CString("Hello, C 0\n") // C.CString 会在C的内存空间申请一个C语言字符串对象，再将Go字符串拷贝到C字符串
//
//	ret := C.SayHello0(&cblob) // &cblob 取C语言对象cblob的地址
//
//	fmt.Println("ret", ret)
//	fmt.Println("repeat_time", cblob.repeat_time)
//
//	C.free(unsafe.Pointer(cblob.str)) // C.CString 申请的C空间内存不会自动释放，需要显示调用C中的free释放
//}
