package taskpool

import (
	"github.com/Zeb-D/go-util/log"
	"time"
)

// 协程池：减少协程创建、销毁的耗时
// 协程池、队列、工作者

//	ITaskPool 线程池：核心数量、最大数量、业余队列存活时间&时间单元、同步队列&队列满后执行策略
type IGoRoutinePool interface {
	Execute(Runnable) (bool, error)
	IsShutdown() bool
	TryTerminate() (bool, error)
	AwaitTermination(unit time.Duration) (bool, error)
	GetTask(allowCoreThreadTimeOut bool) (Runnable, error) //从同步队列取任务
	AddTask(Runnable) (bool, error)                        //增加任务
	Queue() IGoRoutinePoolQueue
	RemoveWorker(string, bool) (bool, error) //取出任务队列
	AddCompletedTasks(int64) bool            //	增加已经完成的任务数
}

//	ITaskPoolQueue 线程安全 Put Pool
type IGoRoutinePoolQueue interface {
	IsEmpty() bool
	Size() int64
	//	直接存任务，可能会失败
	Add(r Runnable) (bool, error)
	//	直接取任务，可能为空
	Take() (Runnable, error)
	//	尝试一定时间内放任务
	Offer(r Runnable, unit time.Duration) (bool, error)
	//	尝试一定时间取任务
	Poll(unit time.Duration) (Runnable, error)
}

//	协程池满了，队列满了，要怎么拒绝任务呢？
type IGoRoutinePoolRejectedPolicy interface {
	RejectedExecution(Runnable, IGoRoutinePool) (bool, error)
}

type Runnable interface {
	Run() (bool, error)
}

//	GoRoutineWork 执行者
type GoRoutineWork struct {
	core             bool
	allowCoreTimeOut bool //是否允许核心协程超时
	status           int
	name             string //工作者名称
	completedTasks   int64
	keepAliveTime    int64 //秒级时间
	task             Runnable
	pool             IGoRoutinePool
}

func NewWork(task Runnable, name string, core bool, pool IGoRoutinePool) *GoRoutineWork {
	return &GoRoutineWork{
		task: task,
		name: name,
		core: core,
		pool: pool,
	}
}

// Run:核心任务允许不停下来
func (w *GoRoutineWork) Work() (bool, error) {
	// 循环去执行任务
	var task Runnable
	for task = w.task; task != nil && w.status != w.status|Stop; {
		ret, err := task.Run()
		if ret {
			w.completedTasks++
		} else if err != nil {
			log.Info(w.name+" 执行任务出错", log.Any("task", task), log.ErrorField(err))
		}
		task, err = w.pool.GetTask(w.core)
		if err != nil {
			log.Info(w.name+" 获取任务出错", log.Any("isCore", w.core), log.ErrorField(err))
		}
	}
	//	把完成数量加上去
	log.Info(w.name+" worker work end", log.Any("完成任务数", w.completedTasks))
	w.pool.AddCompletedTasks(w.completedTasks)
	//	删除worker
	if _, err := w.pool.RemoveWorker(w.name, w.core); err != nil {
		log.Info("workers remove "+w.name, log.ErrorField(err))
	}

	return true, nil
}
