package pulsar

import (
	"context"
	"fmt"
	"github.com/Zeb-D/go-util/log"
	"os/exec"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/TuyaInc/pulsar-client-go/core/sub"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"gopkg.in/natefinch/lumberjack.v2"
)

func newTestClient() Client {
	log.SetGlobalLog("test", false)

	logrus.SetLevel(logrus.DebugLevel)
	jLoger := &lumberjack.Logger{
		Filename: "logrus.log",
	}
	logrus.SetOutput(jLoger)
	addr := "pulsar://127.0.0.1:6650"
	cfg := ClientConfig{
		PulsarAddr: addr,
	}
	return NewClient(cfg)
}

func TestClient_NewConsumer(t *testing.T) {

	client := newTestClient()
	{
		cfg := ConsumerConfig{
			Topic:        "persistent://test/internal/event",
			Subscription: "ferrari",
			Wg:           &sync.WaitGroup{},
		}
		c, err := client.NewConsumer(cfg)
		assert.Nil(t, err)
		defer c.Close()
	}
	// {
	// 	// 未创建租户
	// 	cfg := ConsumerConfig{
	// 		Topic:        "persistent://test1/internal/event",
	// 		Subscription: "ferrari",
	// 		Wg:           &sync.WaitGroup{},
	// 	}
	// 	_, err := client.NewConsumer(cfg)
	// 	assert.NotNil(t, err, err)
	// }
}

func TestClient_NewProducer(t *testing.T) {

	client := newTestClient()
	{
		cfg := ProducerConfig{
			Topic: "persistent://test/internal/event-partition-0",
		}
		_, err := client.NewProducer(cfg)
		assert.Nil(t, err)
	}

	{
		cfg := ProducerConfig{
			Topic: "persistent://test1/internal/event-partition-0",
		}
		_, err := client.NewProducer(cfg)
		assert.NotNil(t, err, err)
	}
}

func TestProducerRetry(t *testing.T) {
	client := newTestClient()
	{
		cfg := ProducerConfig{
			Topic: "persistent://test/internal/event-partition-0",
		}
		p, err := client.NewProducer(cfg)
		assert.Nil(t, err)

		go func() {
			time.Sleep(5 * time.Second)
			restartPulsar()
		}()

		for i := 0; i < 10; i++ {
			payload := []byte("hello-" + fmt.Sprint(i))
			err := p.Publish(context.Background(), payload)
			log.Info("start publish", log.Any("count", i), log.Any("error info", err))
			assert.Nil(t, err)
			time.Sleep(time.Second)
		}
	}
}

func TestConsumerRetry(t *testing.T) {
	client := newTestClient()

	cfg := ConsumerConfig{
		Topic:        "persistent://test/internal/event-partition-0",
		Subscription: "ferrari",
		Wg:           &sync.WaitGroup{},
	}
	cleanTopic(cfg)

	// 消息数量
	count := 20
	// restartPulsarTime 后开始重启pulsar
	restartPulsarTime := 5 * time.Second
	// 每个消息需要消耗的时间
	handlerSpend := 500 * time.Millisecond
	// 所有消息正常情况下需要消耗的时间
	totalSpend := time.Duration(20) * handlerSpend
	// restartPulsarTime必须够小，以保证重启pulsar前消息没有被消费完
	if restartPulsarTime > totalSpend {
		return
	}
	// 所有消息正常情况下需要消耗的时间 加上pulsar重启到恢复正常需要的时间
	timeout := totalSpend + 20*time.Second
	go func() {
		cfg := ProducerConfig{
			Topic: "persistent://test/internal/event-partition-0",
		}
		p, err := client.NewProducer(cfg)
		assert.Nil(t, err)

		for i := 0; i < count; i++ {
			payload := []byte("hello-" + fmt.Sprint(i))
			err := p.Publish(context.Background(), payload)
			log.Debug("start publish", log.Any("count", i), log.Any("error info", err))
			assert.Nil(t, err)
		}
	}()

	c, err := client.NewConsumer(cfg)
	assert.Nil(t, err)

	go func() {
		time.Sleep(restartPulsarTime)
		restartPulsar()
	}()

	h := &HelloHandler{Sleep: handlerSpend}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	c.ReceiveAndHandle(ctx, h)

	c.Close()
	cancel()
	assert.Equal(t, int32(count), h.Count)
}

func cleanTopic(cfg ConsumerConfig) {
	client := newTestClient()
	log.Debug("cleanTopic start")
	c, _ := client.NewConsumer(cfg)
	h := &HelloHandler{}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	c.ReceiveAndHandle(ctx, h)
	c.Close()
	cancel()
	log.Debug("cleanTopic done")
}

func TestCancelConsumer(t *testing.T) {
	client := newTestClient()

	cfg := ConsumerConfig{
		Topic:        "persistent://test/internal/event-partition-0",
		Subscription: "ferrari",
		Wg:           &sync.WaitGroup{},
	}
	c, err := client.NewConsumer(cfg)
	go c.(*consumerImpl).Cronflow()
	assert.Nil(t, err)

	h := &HelloHandler{}

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(10 * time.Second)
		c.Close()
		cancel()
	}()
	c.ReceiveAndHandle(ctx, h)
}

func TestUnactive(t *testing.T) {

}

type HelloHandler struct {
	Sleep time.Duration
	Count int32
	C     *sub.Consumer
	Flag  bool
}

func (h *HelloHandler) HandlePayload(ctx context.Context, msg *Message, payload []byte) error {
	log.Debug("helloHandler received message",
		log.Any("payload", string(payload)),
		log.Any("uuid", msg.Msg.GetMessageId()),
	)
	atomic.AddInt32(&h.Count, 1)
	time.Sleep(h.Sleep)
	return nil
}

func restartPulsar() {
	log.Debug("restartPulsar")
	cmd := exec.Command("docker", "restart", "pulsar-ferrari-test")
	err := cmd.Run()
	if err != nil {
		log.Error("restart docker failed", log.ErrorField(err))
	}
}

func Test_client_NewProducers(t *testing.T) {

	{
		c := newTestClient()
		cfg := ProducersConfig{
			Topic: "persistent://test/internal-1/event-partition-0",
		}
		_, err := c.NewProducers(cfg)
		assert.NotNil(t, err)
	}

	{
		c := newTestClient()
		cfg := ProducersConfig{
			Topic: "persistent://test/internal/event-partition-0",
		}
		ps, err := c.NewProducers(cfg)
		assert.Nil(t, err)
		assert.Equal(t, ps.Partition(), 0)
	}

	{
		c := newTestClient()
		cfg := ProducersConfig{
			Topic: "persistent://test/internal/event",
		}
		ps, err := c.NewProducers(cfg)
		assert.Nil(t, err)
		assert.Equal(t, ps.Partition(), 61)
	}
}
