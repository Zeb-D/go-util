package common

import (
	"fmt"
	"log"
	"sync"
	"time"
)

//守护协程：Supervisor
//当进程中某个goroutine发生panic时，在go的设计理念中，对此是panic零容忍
//在现实中，有可能NPE、数组越界，特别是net/http对每个request的处理
//Supervisor 在处理业务时候，采用柔和的手段，recover捕获，设置重试次数去重启goroutine

type panicInfo struct {
	//异常信息
	key         string
	work        func()
	recoverInfo interface{}
}

func (p *panicInfo) String() string {
	return fmt.Sprintf("panicInfo:{key: %s, recoverInfo:%v}", p.key, p.recoverInfo)
}

type Supervisor struct {
	panicInfoChan  chan panicInfo
	panicPool      map[string]uint
	mux            *sync.Mutex
	maxRestartTime uint
}

func (s *Supervisor) Run(f func(), info string) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				fmt.Printf("[err:%s,info:%s \n]", err, info)
				s.panicInfoChan <- panicInfo{
					key:         info,
					work:        f,
					recoverInfo: err,
				}
			}
		}()
		f()
	}()
}

//MakeSupervisor create goroutine supervisor that proper panic
func MakeSupervisor(maxRestartTime uint) Supervisor {
	s := Supervisor{
		panicInfoChan:  make(chan panicInfo),
		panicPool:      make(map[string]uint), //统计重试次数
		mux:            &sync.Mutex{},
		maxRestartTime: maxRestartTime, //最大重试次数
	}
	//跑一个协程去做恢复操作
	go s.supervisor()
	return s
}

//守护操作: 同一个引用不同key会共用阻塞
func (s *Supervisor) supervisor() {
	for pi := range s.panicInfoChan {
		i := s.panicPool[pi.key]
		if i >= s.maxRestartTime {
			log.Printf("panic and don't restart, time:%d,%v\n", i, pi)
			return
		}
		s.panicPool[pi.key] = i + 1
		//尝试时间间隔翻倍
		time.Sleep(time.Duration(1<<i) * time.Second)
		log.Printf("panic and restart,time:%d, %v\n", s.panicPool[pi.key], pi)
		//再给次机会去重试
		s.Run(pi.work, pi.key)
	}
}

//不同key 阻塞时间互补影响 todo
func (s *Supervisor) dispatchSupervisor() {

}
