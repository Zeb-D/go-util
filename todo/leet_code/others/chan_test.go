package others

import (
	"fmt"
	"reflect"
	"sync"
	"testing"
	"time"
)

func init() {
	//fmt.Println(common.StartCPUProfile("test.prof"))
}

func TestChan(t *testing.T) {

	var ch chan int = make(chan int, 1)
	go func() {
		//ch = make(chan int, 1) //此处打开会输出3
		ch <- 1
		fmt.Println("1")
	}()
	go func(ch chan int) {
		time.Sleep(time.Second)
		<-ch
		fmt.Println("2")
	}(ch)
	//c := time.Tick(1 * time.Second)
	//for range c {
	//	//至少输出2
	//	fmt.Printf("#goroutines: %d\n", runtime.NumGoroutine())
	//}
}

//读已经关闭的 chan 能一直读到东⻄，但是读到的内容根据通道内关闭前是否有元素 而不同。
//如果 chan 关闭前，buffer 内有元素还未读 , 会正确读到 chan 内的值，且返回的第二 个 bool 值(是否读成功)为 true。
//如果 chan 关闭前，buff er 内有元素已经被读完，chan 内无值，接下来所有接收的值 都会非阻塞直接成功，返回 channel 元素的零值，但是第二个 bool 值一直为 false。
func TestReadCloseChan(t *testing.T) {
	ch := make(chan int, 1)
	ch <- 1
	close(ch) //先放一个，关闭，缓冲还存在的
	//ch <- 1 //会报错，关闭的chan
	i, ok := <-ch
	fmt.Println(i, ok) //关闭的chan，读的一直是0
	i, ok = <-ch
	fmt.Println(i, ok)
}

func TestChan2(t *testing.T) {
	var ch chan int // 未初始化，还是nil，未初始化，写和读都会阻塞
	fmt.Println(ch)
	var count int
	go func() {
		ch <- 1
	}()
	go func() {
		count++
		close(ch) //nil 不能close
	}()
	<-ch //nil 可能没有缓冲
	fmt.Println(count)
}

func TestChan3(t *testing.T) {
	i := make(chan interface{}, 1)
	s := make(chan struct{}, 1)
	i <- nil
	i1, ok := <-i
	fmt.Println(i1, "-1->ok,", ok)
	close(i)
	i1, ok = <-i
	fmt.Println(i1, "-2->ok,", ok)
	//i <- nil run error，at close chan
	var a chan interface{}
	fmt.Println("nil chan:", a == nil)
	// a <- nil run nil chan write，will be block
	//s <- nil compile error
	s <- struct{}{}
	s1, ok := <-s
	fmt.Println(s1, "-3->ok,", ok)
	//对s1进行反射
	tv := reflect.TypeOf(s1)
	fmt.Println(tv.Size())
}

func TestSyncMap(t *testing.T) {
	var m sync.Map
	m.LoadOrStore("a", 1)
	m.Delete("a")
	fmt.Println(m)
}

var cc = make(chan int)
var aa int

func ff() {
	aa = 1
	<-cc
}

func TestHappenBefore(t *testing.T) {
	go ff()
	cc <- 0
	print(aa)
}
