package todo

import (
	"testing"
	"time"
)

func TestTimer(t *testing.T) {
	ch := make(chan int)
	tt := time.NewTimer(5 * time.Second)
	select {
	case <-ch:
		// 做相关的业务
	case <-tt.C:
		// 超时了，做超时处理
	}
	tt.Stop() //显示释放资源
}
