package todo

import (
	"context"
	"fmt"
	"strconv"
	"testing"
	"time"
)

func ProcessChannelMessages(ctx context.Context, in <-chan string) {
	for {
		start := time.Now()
		select { //对case都是随机执行，所以他们的case 值都会创建
		case s, ok := <-in:
			if !ok {
				return
			}
			fmt.Println(s)
			// handle `s`
			// 在计时器触发之前，垃圾收集器不会回收 Timer。
			// 因为在上面一进入select这里，time.After 就已经分配内存了，但是上面的ruturn分支
		case <-time.After(5 * time.Minute):
			fmt.Println(start)
		case <-ctx.Done():
			return
		}
	}
}

func ProcessChannelMessagesV2(ctx context.Context, in <-chan string) {
	idleDuration := 5 * time.Minute
	idleDelay := time.NewTimer(idleDuration)
	defer idleDelay.Stop()
	for {
		idleDelay.Reset(idleDuration)
		select {
		case s, ok := <-in:
			if !ok {
				return
			}
			fmt.Println(s)
			// handle `s`
			// 将内存消耗减少 20 倍
		case <-idleDelay.C:
			fmt.Println(idleDelay)
		case <-ctx.Done():
			return
		}
	}
}

type u struct {
	a []*string
}

var u1 = &u{}

func TestTime(t *testing.T) {
	fmt.Println(len(u1.a))
	go Init()
	fmt.Println(len(u1.a))

	in := len(u1.a) > 0
	for !in {
		select {
		case <-time.After(1 * time.Second):
			println(len(u1.a))
			in = len(u1.a) > 0
		}
	}
	fmt.Println(in)
	fmt.Println(u1.a)
}

func Init() {
	time.Sleep(3 * time.Second)
	u1.a = make([]*string, 2)
	u1.a[0] = makeStr()
	u1.a[1] = makeStr()
}

func makeStr() *string {
	str := strconv.Itoa(1)
	return &str
}

func TestForTime(t *testing.T) {
	var sum int = 0
	for i := 0; i < 100; i += 2 {
		for j := 0; j < i; j++ {
			sum++
		}
	}
	fmt.Println(sum)
}
