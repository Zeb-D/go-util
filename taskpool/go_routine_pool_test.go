package taskpool

import (
	"fmt"
	"github.com/Zeb-D/go-util/log"
	"strconv"
	"sync"
	"testing"
	"time"
)

func init() {
	log.SetGlobalLog("taskpool", true)
}

//	执行任务总数
const TestSize = 100000

func TestGoRoutinePool(t *testing.T) {
	q := NewChanTaskQueue(10000)
	p := NewGoRoutinePool(2000, 4000, 3*time.Second, q, &CallerRunsPolicy{})
	var wg sync.WaitGroup
	wg.Add(TestSize)
	start := time.Now()
	log.Info("开始任务")
	fmt.Println("开始任务？")
	for i := 0; i < TestSize; i++ {
		ret, err := p.Execute(&testR{name: strconv.Itoa(i)})
		if err != nil {
			log.Info("p.Execute", log.Any("ret", ret), log.ErrorField(err))
		}
		if i%15000 == 0 {
			time.Sleep(1 * time.Second)
		}
		if p.queue.Size() == 10000 && p.largestPoolSize.Load() >= 3900 {
			time.Sleep(1 * time.Second)
		}
		wg.Done()
	}
	log.Info("start")
	wg.Wait()
	log.Info("end")
	if ret, err := p.AwaitTermination(3 * time.Second); err != nil {
		log.Info("终止条件失败", log.Any("ret", ret), log.ErrorField(err))
	}
	defer log.Info("完成任务数量", log.Any("", p.completedTaskCount.Load()))
	log.Info("总耗时：", log.Any("", time.Now().Sub(start)))
}

func TestGoRoutinePool2(t *testing.T) {
	println(TerminateNoCompleted)
	println(TerminateCompleted)
	println(StopNoCompleted)
	println(StopCompleted)
}

func TestGoRoutinePool_AddWork(t *testing.T) {
	q := NewChanTaskQueue(1)
	p := NewGoRoutinePool(1, 3, 3*time.Second, q, defaultPolicy)
	p.Execute(&testR{})
	p.Execute(&testR{})
	p.Execute(&testR{})
	p.Execute(&testR{})
	p.AwaitTermination(2 * time.Second)
}

func Test(t *testing.T) {
	unit := 0 * time.Second
	println(unit == 0*time.Nanosecond)
}
