package common

import (
	"testing"
	"time"
)

func TestSupervisorRun(t *testing.T) {
	g := MakeSupervisor(3)
	f := func() {
		TimeToDuration(time.Now(), time.Duration(123))
	}
	//一个协程去run,和key无关会阻塞
	g.Run(f, "3")
	g.Run(f, "33")
	time.Sleep(time.Second * 8)
}
