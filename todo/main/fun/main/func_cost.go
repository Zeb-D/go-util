package main

//统计下golang 方法调用的耗时
//go build -gcflags="-m -l" func_cost.go
func main() {
	for i := 0; i < 100000000; i++ {
		a(i)
	}
}

func a(a int) int {
	return 2 * a
}
