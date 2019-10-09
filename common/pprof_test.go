package common

import (
	"fmt"
	"runtime/pprof"
	"testing"
	"time"
)

func TestStartCPUProfile(t *testing.T) {
	StartCPUProfile("test.prof")
	defer pprof.StopCPUProfile()
	for a := 1; a < 100; a++ {
		time.Sleep(100 * time.Millisecond)
		fmt.Println(a)
	}
	//执行
	//go test -v -run TestStartCPUProfile
	//会在当前路径下生成cpu.prof 文件,然后执行
	//go tool pprof test.prof
}
