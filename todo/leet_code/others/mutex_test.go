package others

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

var mu sync.Mutex
var chain string

func TestMutex(t *testing.T) {
	chain = "main"
	AA()
	fmt.Println(chain)
}
func AA() {
	mu.Lock()
	defer mu.Unlock()
	chain = chain + " --> AA"
	BB()
}
func BB() {
	chain = chain + " --> BB"
	CC()
}
func CC() {
	mu.Lock()
	defer mu.Unlock()
	chain = chain + " --> CC"
}

var rwmu sync.RWMutex
var count int

func TestRWMutex(t *testing.T) {
	go A()
	time.Sleep(2 * time.Second)
	fmt.Println("t1")
	mu.Lock()
	fmt.Println("t2")
	defer mu.Unlock()
	fmt.Println("t3")
	count++
	fmt.Println(count)
}
func A() {
	rwmu.RLock()
	defer rwmu.RUnlock()
	B()
	fmt.Println("a")
}
func B() {
	time.Sleep(5 * time.Second)
	C()
	fmt.Println("b")
}
func C() {
	rwmu.RLock()
	defer rwmu.RUnlock()
	fmt.Println("c")
}

type MyMutex struct {
	count int
	sync.Mutex
}

func TestMutex1(t *testing.T) {
	var mu MyMutex
	mu.Lock()
	var mu2 = mu
	mu.count++
	mu.Unlock()
	fmt.Printf("%+v,%+v \n", mu, mu2)
	mu2.Lock() //会发生死锁 fatal error: all goroutines are asleep - deadlock!
	mu2.count++
	mu2.Unlock()
	fmt.Println(mu.count, mu2.count)
}
