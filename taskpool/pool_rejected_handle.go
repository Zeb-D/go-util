package taskpool

import (
	"errors"
)

var AbortPolicyError = errors.New("abort reject execution")
var CallerRunsPolicyError = errors.New("caller runs policy reject execution")
var DiscardPolicyPolicyError = errors.New("discard policy reject execution")
var DiscardOldestPolicyError = errors.New("discard oldest policy reject execution")

//	默认策略为直接拒绝
var defaultPolicy IGoRoutinePoolRejectedPolicy = &AbortPolicy{}

type AbortPolicy struct {
}

func (p *AbortPolicy) RejectedExecution(Runnable, IGoRoutinePool) (bool, error) {
	return false, AbortPolicyError
}

type CallerRunsPolicy struct {
}

func (p *CallerRunsPolicy) RejectedExecution(r Runnable, pool IGoRoutinePool) (bool, error) {
	if !pool.IsShutdown() {
		return r.Run()
	}
	return false, CallerRunsPolicyError
}

type NewGoRoutineRunsPolicy struct {
}

func (p *NewGoRoutineRunsPolicy) RejectedExecution(r Runnable, pool IGoRoutinePool) (bool, error) {
	go r.Run()
	return true, nil
}

//	DiscardPolicy:任务直接丢弃
type DiscardPolicy struct {
}

func (p *DiscardPolicy) RejectedExecution(r Runnable, pool IGoRoutinePool) (bool, error) {
	return false, DiscardPolicyPolicyError
}

// DiscardOldestPolicy:把队列最老的丢了，再执行当前的
type DiscardOldestPolicy struct {
}

func (p *DiscardOldestPolicy) RejectedExecution(r Runnable, pool IGoRoutinePool) (bool, error) {
	pool.Queue().Take() //队列老任务直接丢弃
	pool.Execute(r)
	return false, DiscardOldestPolicyError
}
