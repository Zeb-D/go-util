package interact

import (
	"container/list"
	"flag"
	"fmt"
	"runtime"
	"testing"
	"time"
)

var numCores = flag.Int("n", 2, "number of CPU cores to use")

func TestGOMaxProcs(t *testing.T) {
	flag.Parse()
	fmt.Println(*numCores)
	runtime.GOMAXPROCS(*numCores)

	fmt.Println("In main()")
	go longWait()
	go shortWait()
	fmt.Println("About to sleep in main()")
	// sleep works with a Duration in nanoseconds (ns) !
	time.Sleep(10 * 1e9)
	fmt.Println("At the end of main()")
}

func longWait() {
	fmt.Println("Beginning longWait()")
	time.Sleep(5 * 1e9) // sleep for 5 seconds
	fmt.Println("End of longWait()")
}

func shortWait() {
	fmt.Println("Beginning shortWait()")
	time.Sleep(2 * 1e9) // sleep for 2 seconds
	fmt.Println("End of shortWait()")
}

func TestChan(t *testing.T) {
	ch := make(chan string)

	go sendData(ch)
	//go getData(ch)
	fmt.Println(<-ch)
	time.Sleep(1e9)

	out := make(chan int)
	go func(chan int) {
		o, open := <-out
		if open {
			fmt.Println(o)
		}
	}(out)
	out <- 2 //这里没协程进入，会死锁，必须先有【消费者】
	time.Sleep(2e9)
	close(out)
}

func TestIter(t *testing.T) {
	l := &l{*list.New()}
	l.PushFront("aaa")
	l.PushBack("bbb")
	go func() {
		for ch := range l.Iter() {
			fmt.Println(ch)
		}
	}()
	l.PushFront("ccc")
}

type l struct {
	list.List //匿名字段，相当于继承
}

func (c l) Iter() <-chan interface{} {
	ch := make(chan interface{})
	go func() {
		for next := c.Front(); next.Value != nil; next = next.Next() {
			ch <- next.Value
		}
		return
	}()
	return ch
}

func sendData(ch chan string) {
	ch <- "Washington"
	fmt.Println(11)
	ch <- "Tripoli"
	fmt.Println(22)
	ch <- "London"
	ch <- "Beijing"
	ch <- "Tokio"
}

func getData(ch chan string) {
	var input string
	time.Sleep(1e9)
	for {
		input = <-ch
		fmt.Printf("%s \n", input)
	}
}

var values = [5]int{10, 11, 12, 13, 14}

func TestRoutine(t *testing.T) {
	// 版本A:
	for ix := range values { // ix是索引值
		func() {
			fmt.Print(ix, " ")
		}() // 调用闭包打印每个索引值
	}
	fmt.Println()
	// 版本B: 和A版本类似，但是通过调用闭包作为一个协程
	for ix := range values {
		go func() {
			fmt.Print(ix, " ")
		}()
	}
	fmt.Println()
	time.Sleep(5e9)
	// 版本C: 正确的处理方式
	for ix := range values {
		go func(ix interface{}) {
			fmt.Print(ix, " ")
		}(ix)
	}
	fmt.Println()
	time.Sleep(5e9)
	// 版本D: 输出值:
	for ix := range values {
		val := values[ix]
		go func() {
			fmt.Print(val, " ")
		}()
	}
	time.Sleep(1e9)
}
