package common

import (
	"context"
	"fmt"
	"github.com/magiconair/properties/assert"
	"log"
	"sync"
	"testing"
	"time"
)

func TestContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	ctx.Value("hello")
	cancel()
	b := make(chan bool, 1)
	done := <-ctx.Done()
	fmt.Println("done:", done)
	b <- true
	a := <-b
	assert.Equal(t, true, a)
}

func TestNewContext(t *testing.T) {
	ctx, cancel := NewContext()
	fmt.Sprintf("ctx:%s,cancel:%p \n", ctx, &cancel)

	ch := make(chan bool, 1)
	go func(ctx2 context.Context, ok chan bool) {
		for {
			time.Sleep(1 * time.Second)
			select {
			case <-ctx.Done():
				log.Printf("done")
				ch <- true
			default:
				log.Printf("work")
			}
		}
	}(ctx, ch)
	//cancel()
	time.Sleep(5 * time.Second)
	cancel()
	assert.Equal(t, true, <-ch)
}

func TestNewContextWithTimeout(t *testing.T) {
	//默认超时2s
	ctx, cancel := NewContextWithTimeout(3)
	defer cancel()
	ok := make(chan bool, 1)
	select {
	case <-ctx.Done():
		ok <- true
	case <-time.After(2 * time.Second): //2秒后进入
		ok <- false
	}
	assert.Equal(t, false, <-ok)
}

func TestWait(t *testing.T) {
	wg := sync.WaitGroup{}
	wg.Add(1) //3个信号量
	go func() {
		defer wg.Done() //释放一个
	}()
	fmt.Println("lalallalalal")
	wg.Wait() //当前是否 没有信号量
	fmt.Println("mamamammama")
}