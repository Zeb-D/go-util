

## 协程池

无论协程创建、销毁都比较耗资源，还是同时存在大量的协程会让协程M:N模型调度器更累

### 功能

核心协程：作为核心，当提交任务到协程池中，核心协程数足够则直接创建协程，期间直到服务停止；

任务队列(同步队列)：当核心协程数达到后，不会立即启用业余协程，而是将任务放到任务队列中；

业余协程：业余协程在无任务情况下，允许存活一段时间；当队列满了之后之间启用业务协程；

拒绝策略：当任务队列、业余协程忙中，将采用某种拒绝策略；

停止协程池：慢性通知--先通知池停止提交任务，所有协程后台将所有任务处理完；

主动通知--停止提交任务，等待一段时间后台所有协程处理所有任务，再判断是否还有协程在跑；

终止协程池：把所有协程当前手上的任务都处理好，直接退出；



### 使用

详情请见go_routine_pool_test.go

```go
func init() {
   log.SetGlobalLog("taskpool", true)
}

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

type testR struct {
	name string
}

func (t *testR) Run() (bool, error) {
	time.Sleep(1 * time.Second)
	log.Info(" ->", log.Any("name", t.name))
	return true, nil
}
```



### 设计

GoRoutine Worker

```go
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
```

面向接口编程，提供扩展空间；

```go
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
```

队列主要存放Runnable实现类

```go
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
```

```go
//	协程池满了，队列满了，要怎么拒绝任务呢？
type IGoRoutinePoolRejectedPolicy interface {
	RejectedExecution(Runnable, IGoRoutinePool) (bool, error)
}
```

```go
type Runnable interface {
	Run() (bool, error)
}
```

