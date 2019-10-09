package mq


import (
	"context"
	"log"
	"os"

	"github.com/segmentio/kafka-go"
)

type Message = kafka.Message

type KafkaHandler interface {
	Handle(ctx context.Context, msg *Message) error
}

type ConsumerConfig struct {
	Brokers []string
	Topic   string
	GroupID string
	Debug   bool
}

func NewConsumer(cfg ConsumerConfig) Consumer {
	kcfg := kafka.ReaderConfig{
		Brokers:  cfg.Brokers,
		GroupID:  cfg.GroupID,
		Topic:    cfg.Topic,
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB

	}
	if cfg.Debug {
		kcfg.Logger = log.New(os.Stdout, "[Sarama] ", log.LstdFlags)
	}
	r := kafka.NewReader(kcfg)

	c := &consumer{
		r: r,
	}
	return c
}

type Consumer interface {
	ReceiveAndHandle(h KafkaHandler) error
}

type consumer struct {
	r *kafka.Reader
}

func (c *consumer) ReceiveAndHandle(h KafkaHandler) error {
	defer c.r.Close()
	for {
		m, err := c.r.ReadMessage(context.Background())
		if err != nil {
			return err
		}
		h.Handle(context.Background(), &m)
	}
}

type Producer interface {
	Publish(ctx context.Context, key, payload []byte) error
}

type ProducerConfig struct {
	Brokers []string
	Topic   string
	Debug   bool
}

type producer struct {
	w *kafka.Writer
}

func (p *producer) Publish(ctx context.Context, key, payload []byte) error {
	return p.w.WriteMessages(ctx,
		kafka.Message{
			Key:   key,
			Value: payload,
		},
	)
}

func NewProducer(cfg ProducerConfig) Producer {
	kcfg := kafka.WriterConfig{
		Brokers: cfg.Brokers,
		Topic:   cfg.Topic,
	}
	if cfg.Debug {
		kcfg.Logger = log.New(os.Stdout, "[Sarama] ", log.LstdFlags)
	}
	w := kafka.NewWriter(kcfg)
	p := &producer{
		w: w,
	}
	return p
}


