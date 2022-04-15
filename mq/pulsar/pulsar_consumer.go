package pulsar

import (
	"context"
	"fmt"
	"github.com/Zeb-D/go-util/log"
	"sync"
	"time"

	"go.uber.org/zap"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/TuyaInc/pulsar-client-go/core/manage"
	"github.com/TuyaInc/pulsar-client-go/core/msg"
)

type Message = msg.Message

type Consumer interface {
	ReceiveAndHandle(ctx context.Context, handler PayloadHandler)
	Close() error
}

type PayloadHandler interface {
	HandlePayload(ctx context.Context, msg *Message, payload []byte) error
}

type ConsumerConfig struct {
	Topic        string
	Subscription string
	Wg           *sync.WaitGroup
	// 一般不需要设置，默认为128。模拟消息堆积时，可以设置小一些
	QueueSize int
}

type consumerImpl struct {
	mcsm *manage.ManagedConsumer
	wg   *sync.WaitGroup
	// 因为存在定时flow这个机制，所以有可能出现flow过来的消息进入overflow变量
	// 接下来因为overflow数量大于1，会触发逻辑：向pulsar发送CommandRedeliverUnacknowledgedMessages
	// 注：broker收到CommandRedeliverUnacknowledgedMessages后， failover模式（非独占模式）下，会忽略messagedIDs
	// 由于以上原因，消息可能会乱序以及被重复投递
	// 此时，需要过滤掉重复投递的消息
	acked      map[string]int
	closed     chan struct{}
	flowPeriod time.Duration
}

func (c *consumerImpl) Close() error {
	err := c.mcsm.Close(context.Background())
	if err != nil {
		return err
	}
	close(c.closed)
	return err
}

func (c *consumerImpl) ReceiveAsync(ctx context.Context, queue chan Message) {
	go func() {
		err := c.mcsm.ReceiveAsync(ctx, queue)
		if err != nil {
			log.Error("ReceiveAsync failed", log.ErrorField(err))
		}
	}()
}

// Cronflow 对于单个无分区的topic，初始化后需要显式调用Cronflow，以保证不会消息堆积
// Cronflow不放在ReceiveAsync,因为为每个partition分配一个ticker，代价有些大
func (c *consumerImpl) Cronflow() {
	p := c.flowPeriod
	if p == 0 {
		p = time.Second * 60
	}
	tk := time.NewTicker(p)
	for {
		select {
		case <-tk.C:
			csm := c.mcsm.Consumer(context.Background())
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
					err := csm.Flow(10)
					if err != nil {
						log.Error("flow failed", log.ErrorField(err), log.String("topic", csm.Topic))
					}
				}
			}
		case <-c.closed:
			tk.Stop()
			return
		}
	}
}

// Histogram类型指标，bucket代表duration的分布区间
var MessageHandleDuration = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name:    "pulsar_message_request_duration_seconds",
		Help:    "pulsar message request duration distribution",
		Buckets: []float64{0.005, 0.01, 0.015, 0.02, 0.03, 1},
	},
	[]string{"topic"},
)

func init() {
	prometheus.MustRegister(MessageHandleDuration)
}

func tostring(i int) string {
	list := make([]int, 40)
	for j := 0; j < 40; j++ {
		list[j] = i
	}
	return fmt.Sprint(list)
}

func (c *consumerImpl) ReceiveAndHandle(ctx context.Context, handler PayloadHandler) {
	queue := make(chan Message, 8)
	go c.ReceiveAsync(ctx, queue)
	for {
		select {
		case m := <-queue:
			resp := make(chan struct{})
			go func() {
				now := time.Now()
				c.Handler(ctx, handler, &m, resp)
				MessageHandleDuration.With(prometheus.Labels{"topic": m.Topic}).Observe(time.Since(now).Seconds())
			}()
			select {
			case <-resp:
			case <-time.After(1 * time.Second):
				fields := make([]zap.Field, 0, 10)
				fields = append(fields, log.String("topic", m.Topic))
				fields = append(fields, log.Any("uuid", m.Msg.GetMessageId()))
				log.Warn("force ack", fields...)
			}
		case <-ctx.Done():
			return
		}
	}
}

func (c *consumerImpl) Handler(ctx context.Context, handler PayloadHandler, m *Message, resp chan struct{}) {
	c.wg.Add(1)
	defer func() {
		c.wg.Done()
		close(resp)
	}()

	var list []*msg.SingleMessage
	var err error
	num := m.Meta.GetNumMessagesInBatch()
	if num > 0 && m.Meta.NumMessagesInBatch != nil {
		list, err = msg.DecodeBatchMessage(m)
		if err != nil {
			log.Error("DecodeBatchMessage failed", log.ErrorField(err))
			return
		}
	}

	if c.mcsm.Unactive() {
		log.Warn("unused msg because of consumer is unactivated", log.Any("payload", string(m.Payload)))
		return
	}

	idCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	if c.mcsm.ConsumerID(idCtx) != m.ConsumerID {
		log.Warn("unused msg because of different ConsumerID", log.Any("payload", string(m.Payload)))
		return
	}
	cancel()

	if len(list) == 0 {
		err = handler.HandlePayload(ctx, m, m.Payload)
	} else {
		for i := 0; i < len(list); i++ {
			err = handler.HandlePayload(ctx, m, list[i].SinglePayload)
			if err != nil {
				break
			}
		}
	}
	if err != nil {
		log.Error("handle message failed,dont ack", log.ErrorField(err),
			log.String("topic", m.Topic),
		)
	}

	err = c.mcsm.Ack(ctx, *m)
	if err != nil {
		log.Error("ack failed", log.ErrorField(err))
	}
}
