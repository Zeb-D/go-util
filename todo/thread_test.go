package todo

import (
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"runtime/debug"
	"strconv"
	"sync"
	"testing"
	"time"
)

// 从这次运行可以看出，限制可用的物理线程为 10 个，
// 其中系统占用了 8 个物理线程，user-level 可运行 2 个线程，开启第 3 个线程时就崩溃了。
// 运行结果在不同的 go 版本是不同的，比如 Go1.8 有时候启动 4 到 5 个 goroutine 就会崩溃。
func TestThreadLimit(t *testing.T) {
	nv := 10
	ov := debug.SetMaxThreads(nv)
	fmt.Println(fmt.Sprintf("Change max threads %d=>%d", ov, nv))

	var wg sync.WaitGroup
	c := make(chan bool, 0)
	for i := 0; i < 10; i++ {
		fmt.Println(fmt.Sprintf("Start goroutine #%v", i))

		wg.Add(1)
		go func() {
			c <- true
			defer wg.Done()
			runtime.LockOSThread()
			time.Sleep(10 * time.Second)
			fmt.Println("Goroutine quit")
		}()

		<-c
		fmt.Println(fmt.Sprintf("Start goroutine #%v ok", i))
	}

	fmt.Println("Wait for all goroutines about 10s...")
	wg.Wait()

	fmt.Println("All goroutines done")
}

// 如何避免程序超过线程限制被干掉？
// 一般可能阻塞在 system call，那么什么时候会阻塞？还有，GOMAXPROCS 又有什么作用呢？
func TestThreadPanic(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("main recover is", r)
		}
	}()

	nv := 10
	ov := debug.SetMaxThreads(nv)
	fmt.Println(fmt.Sprintf("Change max threads %d=>%d", ov, nv))

	var wg sync.WaitGroup
	c := make(chan bool, 0)
	for i := 0; i < 10; i++ {
		fmt.Println(fmt.Sprintf("Start goroutine #%v", i))

		wg.Add(1)
		go func() {
			c <- true

			defer func() {
				if r := recover(); r != nil {
					fmt.Println("main recover is", r)
				}
			}()

			defer wg.Done()
			runtime.LockOSThread()
			time.Sleep(10 * time.Second)
			fmt.Println("Goroutine quit")
		}()

		<-c
		fmt.Println(fmt.Sprintf("Start goroutine #%v ok", i))
	}

	fmt.Println("Wait for all goroutines about 10s...")
	wg.Wait()

	fmt.Println("All goroutines done")
}

// 虽然 GOMAXPROCS 设置为 1，实际上创建了 12 个物理线程。有大量的时间是在 sys 上面，也就是 system calls。
// time go run t.go
//real    1m44.679s
//user    0m0.230s
//sys    0m53.474s
func TestThreadOk(t *testing.T) {
	runtime.GOMAXPROCS(1)
	data := make([]byte, 128*1024)

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			for {
				ioutil.WriteFile("testxxx"+strconv.Itoa(n), []byte(data), os.ModePerm)
			}
		}(i)
	}

	wg.Wait()
}
