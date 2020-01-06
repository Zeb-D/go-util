package todo

import (
	"log"
	"sync"
	"testing"
)

const Size = 1024

func init() {
	log.Println("size:", Size)
	log.Println("pool limit:", len(Limit))
}

var P1 = sync.Pool{
	New: func() interface{} {
		return make([]byte, Size)
	},
}

var P2 = sync.Pool{
	New: func() interface{} {
		return make([]byte, Size)
	},
}

var Limit = make(chan struct{}, 10)

// 不使用pool
func Write1(bs []byte) {
	buf := make([]byte, Size)
	num := copy(buf, bs)
	_ = num
}

// 使用pool
func Write2(bs []byte) {
	buf := P1.Get().([]byte)
	num := copy(buf, bs)
	_ = num
	P1.Put(buf)
}

// 使用pool+限制池容量
func Write3(bs []byte) {
	Limit <- struct{}{}
	buf := P2.Get().([]byte)
	num := copy(buf, bs)
	_ = num
	P2.Put(buf)
	<-Limit
}

//go test -v -bench BenchmarkWrite1 -benchmem -run BenchmarkWrite1
// -memprofile mem1 -cpuprofile cpu1 pool_test.go
func BenchmarkWrite1(b *testing.B) {
	bs := make([]byte, 500)
	for i := 0; i < b.N; i++ {
		wg := sync.WaitGroup{}
		for j := 0; j < 1000; j++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				Write1(bs)
			}()
		}
		wg.Wait()
	}
}

//go test -v -bench BenchmarkWrite2 -benchmem -run BenchmarkWrite2 -memprofile mem2 -cpuprofile cpu2 pool_test.go
func BenchmarkWrite2(b *testing.B) {
	bs := make([]byte, 500)
	for i := 0; i < b.N; i++ {
		wg := sync.WaitGroup{}
		for j := 0; j < 1000; j++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				Write2(bs)
			}()
		}
		wg.Wait()
	}
}

//go test -v -bench BenchmarkWrite3 -benchmem -run BenchmarkWrite3 -memprofile mem3 -cpuprofile cpu3 pool_test.go
func BenchmarkWrite3(b *testing.B) {
	bs := make([]byte, 500)
	for i := 0; i < b.N; i++ {
		wg := sync.WaitGroup{}
		for j := 0; j < 1000; j++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				Write3(bs)
			}()
		}
		wg.Wait()
	}
}
