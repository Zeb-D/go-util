package taskpool

import (
	"errors"
	"fmt"
	"github.com/Zeb-D/go-util/log"
	"go.uber.org/atomic"
	"math/rand"
	"sync"
	"time"
)

const (
	Init                 = 0
	Running              = 1
	CorePoolRunning      = 2
	MaxPoolRunning       = 3
	Terminate            = 4
	Stop                 = 5
	TerminateNoCompleted = Terminate<<5 + Terminate
	TerminateCompleted   = Terminate<<6 + Terminate
	StopNoCompleted      = Stop<<5 + Stop
	StopCompleted        = Stop<<6 + Stop
)

var RunnableNil = errors.New("runnable is nil")
var Nil = errors.New("anything is nil")
var ExecuteRunnableError = errors.New("execute runnable error")
var ExecuteAddWorkerError = errors.New("execute add worker error")
var PoolTerminateStatusError = errors.New("pool terminate status error")
var PoolStopStatusError = errors.New("pool stop status error")
var PoolStatusError = errors.New("pool status error")

//	协程池
type goRoutinePool struct {
	//是否允许 核心协程完成任务关闭
	allowCoreThreadTimeOut bool
	//0:初始化，1：启动，2：小于核心协程数在跑，3：大于核心协程数在跑，5：强行停止提交任务，4：停止提交任务
	//6：停止提交任务且任务未完成，7：停止提交任务且任务全部完成
	status             int
	largestPoolSize    *atomic.Int64             //核心协程数量
	corePoolSize       int64                     //核心协程数量
	maxPoolSize        int64                     //核心协程数量
	completedTaskCount *atomic.Int64             //完成的任务数
	keepAliveTime      time.Duration             //闲的协程大概可以空闲多久
	workers            map[string]*GoRoutineWork //工作协程
	queue              IGoRoutinePoolQueue
	rejectedPolicy     IGoRoutinePoolRejectedPolicy
	m                  *sync.Mutex
}

func (p *goRoutinePool) AddWork(firstTask Runnable, core bool) bool {
	// 如果是stop直接强行停止提交任务
	if p.status >= p.status|Stop || firstTask == nil {
		return false
	}
	workerStarted := false
	workerAdded := false
	workName := fmt.Sprintf("work-%v-%d-%d", core, p.largestPoolSize.Load(), rand.Intn(10))
	var w = NewWork(firstTask, workName, core, p)
	// 终止任务	|| (p.status == (p.status|Terminate)
	if p.status < Terminate && firstTask != nil {
		p.m.Lock()
		p.workers[workName] = w
		workerAdded = true
		p.m.Unlock()
	}
	if workerAdded {
		go w.Work()
		workerStarted = true
	}
	//重新计算下工作者大小
	p.largestPoolSize = atomic.NewInt64(int64(len(p.workers)))
	return workerStarted
}

func (p *goRoutinePool) AddCompletedTasks(count int64) bool {
	p.completedTaskCount.Add(count)
	return true
}

func (p *goRoutinePool) Execute(r Runnable) (bool, error) {
	if r == nil {
		return false, RunnableNil
	}
	//	协程池还在跑
	if p.status >= Init && p.status < Terminate {
		// 小于核心协程数，直接创建一个协程去work
		if p.largestPoolSize.Load() < p.corePoolSize {
			if p.AddWork(r, true) {
				return true, nil
			} else {
				return false, PoolStatusError
			}
		}

		//	协程数不小于核心数，那么直接放到队列，队列放不下，则创建业余协程
		ret, _ := p.queue.Add(r)
		if ret {
			return true, nil
		} else {
			//判断协程数量合法？
			if p.largestPoolSize.Load() <= p.maxPoolSize {
				if p.AddWork(r, false) {
					return true, nil
				} else {
					return false, PoolStatusError
				}

			} else {
				//拒绝提交任务
				return p.rejectedPolicy.RejectedExecution(r, p)
			}
		}

	}
	log.Info("还有条件没判断到?",
		log.Any("p.largestPoolSize", p.largestPoolSize.Load()),
		log.Any("p.poolSize", len(p.workers)),
		log.Any("status", p.status))
	return false, ExecuteRunnableError
}

func (p *goRoutinePool) IsShutdown() bool {
	return len(p.workers) == 0 || p.status >= Terminate //	是否被关闭了任务
}

//	TryTerminate 停止提交任务，直到任务全部完成，然后所有协程都会自己完成
func (p *goRoutinePool) TryTerminate() (bool, error) {
	// 停止提交任务
	p.status = Terminate
	//通知工作者，设置终止状态
	for _, work := range p.workers {
		if work == nil {
			continue
		}
		work.status = Terminate
	}

	return true, nil
}

func (p *goRoutinePool) AwaitTermination(unit time.Duration) (bool, error) {
	p.m.Lock()
	defer p.m.Unlock()
	//	参数为0，表示立即停止
	if unit == 0*time.Nanosecond {
		p.status = Stop
		for _, work := range p.workers {
			if work == nil {
				continue
			}
			work.status = Stop
		}
		return true, nil
	}

	p.TryTerminate()
	time.Sleep(unit)
	return p.status == TerminateCompleted, nil
}

func (p *goRoutinePool) RemoveWorker(name string, core bool) (bool, error) {
	if p.workers == nil || len(p.workers) == 0 {
		return false, Nil
	}
	p.m.Lock()
	delete(p.workers, name)
	p.m.Unlock()
	p.largestPoolSize = atomic.NewInt64(int64(len(p.workers)))
	//	计算状态:核心协程都要停下来表示开始结束
	if core && p.largestPoolSize.Load() == 0 {
		if p.status == (p.status | Terminate) {
			p.status = TerminateCompleted
		} else {
			p.status = StopCompleted
		}
	}
	return true, nil
}

//	从同步队列取任务 (bool, error)
func (p *goRoutinePool) GetTask(core bool) (r Runnable, err error) {
	// 这是一个死循环，配合核心协程一直跑
	for {
		//	停止提交任务，队列为空，则直接退出
		if p.status == (p.status|Terminate) && p.queue.IsEmpty() {
			return nil, PoolTerminateStatusError
		}
		// 强行停止任务，则直接退出
		if p.status == (p.status | Stop) {
			return nil, PoolStopStatusError
		}

		// 判断是否要阻塞获取任务，默认核心任务阻塞
		timed := !p.allowCoreThreadTimeOut && core
		r, err = p.queue.Poll(p.keepAliveTime)
		if err != nil {
			log.Info("queue.Poll error", log.ErrorField(err), log.Any("p.queue.size", p.queue.Size()), log.Any("status", p.status))
		}
		// 如果核心协程没有获取到任务则继续
		if r == nil && timed {
			continue
		}
		return r, nil
	}

	return nil, nil
}

//	从同步队列取任务 (bool, error)
func (p *goRoutinePool) AddTask(Runnable) (bool, error) {
	return false, nil
}

func (p *goRoutinePool) Queue() IGoRoutinePoolQueue {
	return p.queue
}

func GoRoutinePool(allowCoreThreadTimeOut bool,
	corePoolSize, maxPoolSize int64,
	keepAliveTime time.Duration,
	queue IGoRoutinePoolQueue,
	policy IGoRoutinePoolRejectedPolicy) *goRoutinePool {
	return &goRoutinePool{
		allowCoreThreadTimeOut: allowCoreThreadTimeOut,
		largestPoolSize:        atomic.NewInt64(0),
		corePoolSize:           corePoolSize,
		maxPoolSize:            maxPoolSize,
		completedTaskCount:     atomic.NewInt64(0),
		keepAliveTime:          keepAliveTime,
		status:                 Init,
		workers:                make(map[string]*GoRoutineWork),
		queue:                  queue,
		rejectedPolicy:         policy,
		m:                      &sync.Mutex{},
	}
}

func NewGoRoutinePool(corePoolSize, maxPoolSize int64, keepAliveTime time.Duration, queue IGoRoutinePoolQueue, rejectedPolicy IGoRoutinePoolRejectedPolicy) *goRoutinePool {
	return GoRoutinePool(false, corePoolSize, maxPoolSize, keepAliveTime, queue, rejectedPolicy)
}
