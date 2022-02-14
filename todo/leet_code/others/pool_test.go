package others

import (
	"bytes"
	"fmt"
	"runtime"
	"sync"
	"testing"
	"time"
)

var pool = sync.Pool{New: func() interface{} { return new(bytes.Buffer) }}

// 这里考察两个知识点，pool + Buffer的grow模式
func TestPool(t *testing.T) {
	go func() {
		for {
			processRequest(1 << 10) // 256MiB
		}
	}()
	for i := 0; i < 1; i++ {
		go func() {
			for {
				processRequest(1 << 10) // 1KiB
			}
		}()
	}
	var stats runtime.MemStats
	for i := 0; ; i++ {
		runtime.ReadMemStats(&stats)
		fmt.Printf("Cycle %d: %dB\n", i, stats.Alloc)
		time.Sleep(time.Second)
		runtime.GC()
	}
}

// Get不到才会new 个*bytes.Buffer，否则使用这个池子的引用
func processRequest(size int) {
	b := pool.Get().(*bytes.Buffer)
	time.Sleep(500 * time.Millisecond)
	b.Grow(size)
	pool.Put(b)
	time.Sleep(1 * time.Millisecond)
}

// 测试Buffer的Grow模式 是否一直在增加
func TestBufferGrow(t *testing.T) {
	b := new(bytes.Buffer)
	for i := 0; i <= 10; i++ {
		go func(i int) {
			time.Sleep(500 * time.Millisecond)
			b.Grow(i << 10)
			fmt.Println(b.Cap())
			time.Sleep(1 * time.Millisecond)
		}(i)
	}

	var stats runtime.MemStats
	for i := 0; ; i++ {
		runtime.ReadMemStats(&stats)
		fmt.Printf("Cycle %d: %dB\n", i, stats.Alloc)
		time.Sleep(time.Second)
		runtime.GC()
	}
}

// 对象池，不够的时候会自己new，只要GC了，会清空
var p = sync.Pool{
	New: func() interface{} {
		return "123"
	},
}

func TestPool2(tt *testing.T) {
	t := p.Get().(string)
	fmt.Println(t)

	p.Put("321")
	t2 := p.Get().(string)
	fmt.Println(t2)
}
