package taskpool

import (
	"errors"
	"fmt"
	"go.uber.org/atomic"
	"sync"
	"time"
)

var QueueFull = errors.New("Queue full")
var QueueEmpty = errors.New("Queue empty")

// 基于chan 实现同步队列
type chanTaskQueue struct {
	size      atomic.Int64
	maxSize   int64
	mux       *sync.Mutex
	workQueue chan Runnable
}

func (q *chanTaskQueue) String() string {
	return fmt.Sprintf("size:%d,maxSize:%d,workQueue:%v", q.size, q.maxSize, q.workQueue)
}

func NewChanTaskQueue(maxQueueSize int64) *chanTaskQueue {
	workQueue := make(chan Runnable, maxQueueSize)
	return &chanTaskQueue{
		maxSize:   maxQueueSize,
		workQueue: workQueue,
		mux:       &sync.Mutex{},
	}
}
func (q *chanTaskQueue) IsEmpty() bool {
	if q.size.Load() != int64(len(q.workQueue)) {
		fmt.Errorf("并发出现了大小不一样size:%d,maxSize:%d，queueSize:%d",
			q.size, q.maxSize, len(q.workQueue))
	}
	return q.size.Load() == 0 && len(q.workQueue) == 0
}

func (q *chanTaskQueue) Size() int64 {
	return int64(len(q.workQueue))
}

func (q *chanTaskQueue) Add(r Runnable) (bool, error) {
	q.mux.Lock()
	defer q.mux.Unlock()
	if q.size.Load() < q.maxSize {
		q.workQueue <- r
		q.size.Inc()
		return true, nil
	}
	return false, QueueFull
}

func (q *chanTaskQueue) Take() (Runnable, error) {
	q.mux.Lock()
	if q.IsEmpty() {
		return nil, QueueEmpty
	}
	r := <-q.workQueue
	q.size.Dec()
	q.mux.Unlock()
	return r, nil
}

func (q *chanTaskQueue) Offer(r Runnable, unit time.Duration) (bool, error) {
	select {
	case <-time.After(unit):
		return false, QueueFull
	case q.workQueue <- r:
		q.size.Inc()
		return true, nil
	}
}

func (q *chanTaskQueue) Poll(unit time.Duration) (Runnable, error) {
	select {
	case <-time.After(unit):
		return nil, QueueEmpty
	case r := <-q.workQueue:
		q.size.Dec()
		return r, nil
	}
}
