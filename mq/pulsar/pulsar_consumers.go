package pulsar

import (
	"context"
	"sync"
	"time"

	"github.com/Zeb-D/go-util/log"
)

type ConsumerList struct {
	list             []*consumerImpl
	FlowPeriodSecond int
	FlowPermit       uint32
	closed           chan struct{}
}

func (l *ConsumerList) ReceiveAndHandle(ctx context.Context, handler PayloadHandler) {
	wg := sync.WaitGroup{}
	for i := 0; i < len(l.list); i++ {
		safe := i
		wg.Add(1)
		go func() {
			l.list[safe].ReceiveAndHandle(ctx, handler)
			wg.Done()
		}()
	}
	go l.CronFlow()
	wg.Wait()
}

func (l *ConsumerList) Close() error {
	wg := sync.WaitGroup{}
	for i := 0; i < len(l.list); i++ {
		safe := i
		wg.Add(1)
		go func() {
			_ = l.list[safe].Close()
			wg.Done()
		}()
	}
	wg.Wait()
	close(l.closed)
	return nil
}

func (l *ConsumerList) CronFlow() {
	if l.FlowPeriodSecond == 0 {
		return
	}
	if l.FlowPermit == 0 {
		return
	}
	tk := time.NewTicker(time.Duration(l.FlowPeriodSecond) * time.Second)
	for {
		select {
		case <-tk.C:
			for i := 0; i < len(l.list); i++ {
				csm := l.list[i].mcsm.Consumer(context.Background())
				if csm != nil {
					if len(csm.Overflow) > 0 {
						log.Info("RedeliverOverflow",
							log.Any("topic", csm.Topic),
							log.Any("num", len(csm.Overflow)),
							log.Any("detail", csm.Overflow),
						)
						_, err := csm.RedeliverOverflow(context.Background())
						if err != nil {
							log.Warn("RedeliverOverflow failed",
								log.Any("topic", csm.Topic),
								log.ErrorField(err),
							)
						}
					}
					if !csm.Unactive && len(csm.Queue) == 0 {
						err := csm.Flow(l.FlowPermit)
						if err != nil {
							log.Error("flow failed", log.ErrorField(err), log.String("topic", csm.Topic))
						}
					}
				}
			}
		case <-l.closed:
			return
		}
	}
}
