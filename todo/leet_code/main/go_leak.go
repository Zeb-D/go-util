package main

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"runtime"
	"time"
)

func toLeak() {
	c := make(chan int)
	go func() {
		<-c
	}()
}

func main() {
	go toLeak()

	go func() {
		_ = http.ListenAndServe("0.0.0.0:8080", nil)
	}()

	c := time.Tick(time.Second)
	for range c {
		fmt.Printf("goroutine [nums]: %d\n", runtime.NumGoroutine())
	}
}

//http://127.0.0.1:8080/debug/pprof/goroutine?debug=1
//http://127.0.0.1:8080/debug/pprof/  在这个页面上，可以看不同的维度信息
//每请求一次，使用go func 处理请求

//一个goroutine启动后没有正常退出，而是直到整个服务结束才退出，这种情况下，goroutine无法释放，内存会飙高，严重可能会导致服务不可用。
//goroutine的退出其实只有以下几种方式可以做到：
//main函数退出
//context通知退出
//goroutine panic退出
//goroutine 正常执行完毕退出
//大多数引起goroutine泄露的原因基本上都是如下情况：
//channel阻塞，导致协程永远没有机会退出
//异常的程序逻辑(比如循环没有退出条件)

//排查方式
//排查:
//go pprof工具
//runtime.NumGoroutine()判断实时协程数
//第三方库

//下面更详细：
//go tool pprof -http=:8001 http://127.0.0.1:8080/debug/pprof/goroutine?debug=1
