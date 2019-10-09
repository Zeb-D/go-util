package mq

import (
	"context"
	"fmt"
	"testing"
)

func TestKafkaConsumer(t *testing.T) {
	cfg := ConsumerConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   "ydkafka",
		GroupID: "testgroup",
		Debug:   true,
	}
	c := NewConsumer(cfg)
	h := &helloHandler{}
	c.ReceiveAndHandle(h)
}

func TestKafkaProducer(t *testing.T) {
	cfg := ProducerConfig{
		Brokers: []string{"localhost:9092"}, //[]string{"172.16.248.136:9092", "172.16.248.37:9092"},
		Topic:   "ydkafka",
	}
	p := NewProducer(cfg)
	p.Publish(context.Background(), nil, []byte(`{"id":1}`))
}

type helloHandler struct {
}

func (h *helloHandler) Handle(ctx context.Context, m *Message) error {
	fmt.Printf("message at topic/partition/offset %v/%v/%v: %s = %s\n", m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))
	return nil
}
