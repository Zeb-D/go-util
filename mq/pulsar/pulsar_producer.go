package pulsar

import (
	"context"
	"github.com/Zeb-D/go-util/log"

	"github.com/TuyaInc/pulsar-client-go/core/pub"
)

type Producer interface {
	Publish(ctx context.Context, payload []byte) error
}

type producer struct {
	p            *pub.Producer
	client       *client
	PublishRetry int
	Topic        string
	Config       ProducerConfig
}

func (p *producer) Publish(ctx context.Context, payload []byte) error {
	var err error
	addr := p.client.config.PulsarAddr

	sCtx, cancel := context.WithTimeout(ctx, p.Config.PublishTimeout)
	_, err = p.p.Send(sCtx, payload)
	cancel()
	if err != nil {
		np, err := p.client.reconnectProducer(p.Config)
		if err != nil {
			return err
		}
		p.p = np
		sCtx, cancel := context.WithTimeout(ctx, p.Config.PublishTimeout)
		_, err = p.p.Send(sCtx, payload)
		cancel()
		if err != nil {
			log.Error("send payload failed after retry",
				log.Any("addr", addr),
				log.Any("topic", p.Topic),
				log.ErrorField(err),
			)
		}
		return err
	}
	return err
}
