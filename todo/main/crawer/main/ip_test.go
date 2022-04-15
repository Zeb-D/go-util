package main

import (
	"fmt"
	"sync/atomic"
	"testing"
)

func TestIpGen(t *testing.T) {
	for i := 0; i < 1000; i++ {
		fmt.Println(genIpAddr())
		//time.Sleep(time.Second)
	}
	var a int32 = 1
	atomic.LoadInt32(&a)

}
