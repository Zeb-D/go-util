package todo

import (
	"fmt"
	"golang.org/x/sys/cpu"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

type Origin struct {
	a uint64
	b uint64
}

var num = 1000 * 1000

func OriginParallel() {
	var v Origin

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		for i := 0; i < num; i++ {
			atomic.AddUint64(&v.a, 1)
		}
		wg.Done()
	}()

	go func() {
		for i := 0; i < num; i++ {
			atomic.AddUint64(&v.b, 1)
		}
		wg.Done()
	}()

	wg.Wait()
	_ = v.a + v.b
}

// cpu line 64kb，otherwise difference from arm
type WithPadding struct {
	a uint64
	_ [56]byte //Padding
	b uint64
	_ [56]byte
}

func WithPaddingParallel() {
	_ := cpu.CacheLinePad{}
	var v WithPadding

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		for i := 0; i < num; i++ {
			atomic.AddUint64(&v.a, 1)
		}
		wg.Done()
	}()

	go func() {
		for i := 0; i < num; i++ {
			atomic.AddUint64(&v.b, 1)
		}
		wg.Done()
	}()

	wg.Wait()
	_ = v.a + v.b
}

// 由于一个 CPU 核在读取一个变量时，以 cache line 的方式将后续的变量也读取进来，缓存在自己这个核的 cache 中，
// 而后续的变量也可能被其他 CPU 核并行缓存。
// 当前面的 CPU 对前面的变量进行写入时，该变量同样是以 cache line 为单位写回内存。
// 此时，在其他核上，尽管缓存的是该变量之后的变量，
// 但是由于没法区分自身变量是否被修改，所以它只能认为自己的缓存失效，重新从内存中读取。
// 这种现象叫做false sharing。
func TestCpuTest(t *testing.T) {
	var b time.Time

	b = time.Now()
	OriginParallel()
	fmt.Printf("OriginParallel. Cost=%+v.\n", time.Now().Sub(b))

	b = time.Now()
	WithPaddingParallel()
	fmt.Printf("WithPaddingParallel. Cost=%+v.\n", time.Now().Sub(b))
}
