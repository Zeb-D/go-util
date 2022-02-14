//cgo_export.go
package main

/*
#include "cgo_export.h"                       // _cgo_export.h由cgo工具动态生成
GoInt32 Add(GoInt32 param1, GoInt32 param2) {       // GoInt32即为cgo在C语言的导出类型
  return param1 + param2;
}

*/
import "C"
import "fmt"

func main() {
	// _Ctype_                      // _Ctype_ 会在cgo预处理阶段触发异常，
	fmt.Println(C.Add(1, 2))
}
