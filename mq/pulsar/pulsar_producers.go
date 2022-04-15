package pulsar

import (
	"context"
	"errors"
	"fmt"
	"github.com/Zeb-D/go-util/log"
	"sync"

	"github.com/TuyaInc/pulsar-client-go/core/manage"
)

type Producers interface {
	Publish(ctx context.Context, payload []byte, partition uint32) error
	Partition() int
}

type ProducersConfig struct {
	Topic string
	// 发送消息，失败重试次数，默认为3
	PublishRetry int
}

func (c *client) NewProducers(cfg ProducersConfig) (Producers, error) {
	ps := &producers{M: sync.Map{}}

	errs := make(chan error, 10)
	go func() {
		for err := range errs {
			log.Error("async errors", log.ErrorField(err))
		}
	}()
	clientConfig := manage.ClientConfig{
		Addr: c.config.PulsarAddr,
		Errs: errs,
	}

	pNum, err := c.GetPartition(cfg.Topic, clientConfig)
	if err != nil {
		return nil, err
	}
	ps.PartitionCount = pNum
	// 有分区
	if pNum > 0 {
		for i := 0; i < pNum; i++ {
			pConfig := ProducerConfig{
				Topic: fmt.Sprintf("%s-partition-%d", cfg.Topic, i),
			}
			p, err := c.NewProducer(pConfig)
			if err != nil {
				return nil, err
			}
			ps.M.Store(i, p)
		}
		return ps, nil
	}
	// 无分区
	pConfig := ProducerConfig{
		Topic: cfg.Topic,
	}
	p, err := c.NewProducer(pConfig)
	if err != nil {
		return nil, err
	}
	ps.P = p
	return ps, nil
}

type producers struct {
	M              sync.Map
	PartitionCount int
	P              Producer
}

var ErrNoPartitionProducer = errors.New("no partition producer")

func (ps *producers) Publish(ctx context.Context, payload []byte, partition uint32) error {
	if ps.PartitionCount > 0 {
		p, ok := ps.M.Load(int(partition))
		if !ok {
			return ErrNoPartitionProducer
		}
		pp := p.(Producer)
		err := pp.Publish(ctx, payload)
		return err
	}
	err := ps.P.Publish(ctx, payload)
	return err
}

func (ps *producers) Partition() int {
	return ps.PartitionCount
}
