package pulsar

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/TuyaInc/pulsar-client-go/core/manage"
	"github.com/TuyaInc/pulsar-client-go/core/pub"
	"github.com/Zeb-D/go-util/common"
	"github.com/Zeb-D/go-util/log"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

const (
	TraceKey                = "trace-key"
	TagPulsarMQTopic        = "pulsar.topic"
	TagPulsarMQPayload      = "pulsar.payload"
	TagPulsarMQBatchPayload = "pulsar.batch.payload"
	TagPulsarMQConsumerMode = "pulsar.consumer.mode"

	// ForTopicTimeout ForTopic最大执行时常
	ForTopicTimeout time.Duration = 5 * time.Second

	DefaultFlowPeriodSecond = 60
	DefaultFlowPermit       = 10
)

type Client interface {
	NewConsumer(config ConsumerConfig) (Consumer, error)
	NewProducer(config ProducerConfig) (Producer, error)
	NewProducers(config ProducersConfig) (Producers, error)
}

type client struct {
	pool   *manage.ClientPool
	config ClientConfig
}

type ProducerConfig struct {
	Topic string
	// 发送消息，失败重试次数，默认为3
	PublishRetry int
	// send 超时时间 默认为5s
	PublishTimeout time.Duration
}

type ClientConfig struct {
	PulsarAddr            string
	Auth                  AuthProvider
	TlsInsecureSkipVerify bool
}

func (c *client) newProducer(config ProducerConfig) (*pub.Producer, error) {
	errs := make(chan error, 10)
	clientCfg := manage.ClientConfig{
		Addr: c.config.PulsarAddr,
		Errs: errs,
	}
	if c.config.Auth != nil {
		clientCfg.AuthData = c.config.Auth.AuthData()
		clientCfg.AuthMethod = c.config.Auth.AuthMethod()
	}
	ctx, cancel := context.WithTimeout(context.Background(), ForTopicTimeout)
	defer cancel()

	topic := config.Topic
	producerName := ""
	mClient, err := c.pool.ForTopic(ctx, clientCfg, topic)
	if err != nil {
		log.Error("clientPool forTopic failed", log.ErrorField(err))
		return nil, err
	}
	ctx, cancel = common.NewContextWithTimeout(ForTopicTimeout)
	defer cancel()
	client, err := mClient.Get(ctx)
	if err != nil {
		log.Error("get pulsar manage client failed", log.ErrorField(err))
		return nil, err
	}

	ctx, cancel = common.NewContextWithTimeout(ForTopicTimeout)
	defer cancel()
	// NewProducer 比较耗时
	p, err := client.NewProducer(ctx, topic, producerName)
	if err != nil {
		log.Error("new pulsar producer for client failed", log.ErrorField(err))
		return nil, err
	}
	p.AddTraceHook(&traceHook{topic: topic})
	return p, nil
}

func (c *client) GetPartition(topic string, config manage.ClientConfig) (int, error) {
	p, err := c.pool.Partitions(context.Background(), config, topic)
	if err != nil {
		log.Error("get partition failed", log.ErrorField(err), log.String("topic", topic))
		return 0, err
	}
	if p.Error != nil {
		log.Error("get partition failed",
			log.String("serve resp", p.GetMessage()),
			log.String("topic", topic),
		)
		return 0, errors.New("get partition failed")
	}
	return int(p.GetPartitions()), nil
}

func (c *client) NewProducer(config ProducerConfig) (Producer, error) {
	p, err := c.newProducer(config)
	if err != nil {
		return nil, err
	}
	// 默认三次重试
	if config.PublishRetry == 0 {
		config.PublishRetry = 3
	}
	if config.PublishTimeout == 0 {
		config.PublishTimeout = time.Second * 5
	}
	return &producer{
		p:            p,
		client:       c,
		Topic:        config.Topic,
		PublishRetry: config.PublishRetry,
		Config:       config,
	}, nil
}

func (c *client) NewConsumer(config ConsumerConfig) (Consumer, error) {
	if config.Wg == nil {
		config.Wg = &sync.WaitGroup{}
	}
	if config.Subscription == "" {
		config.Subscription = subscriptionName(config.Topic)
	}

	errs := make(chan error, 10)
	go func() {
		for err := range errs {
			log.Error("async errors", log.String("topic", config.Topic), log.ErrorField(err))
		}
	}()

	cfg := manage.ConsumerConfig{
		ClientConfig: manage.ClientConfig{
			Addr: c.config.PulsarAddr,
			Errs: errs,
		},

		QueueSize: config.QueueSize, // 测试overflow时，设置为2
		Topic:     config.Topic,
		SubMode:   manage.SubscriptionModeFailover,
		Name:      config.Subscription,
		// NewConsumerTimeout 尽量设置长一点。如果太小，会出现第一次发送Subscribe命令后
		// 由于超过ConsumerTimeout还未得到响应，客户端认为失败（而实际是成功的，只是响应较慢）。
		// 接下来客户端会再次发送Subscribe，第二次虽然发送成功，但是由于已经存在了一个consumer
		// pulsar会认为第二个consumer为unactive状态。此时会导致消息阻塞
		NewConsumerTimeout: time.Minute,
	}
	if c.config.Auth != nil {
		cfg.ClientConfig.AuthData = c.config.Auth.AuthData()
		cfg.ClientConfig.AuthMethod = c.config.Auth.AuthMethod()
	}
	if c.config.TlsInsecureSkipVerify {
		cfg.ClientConfig.TLSConfig = &tls.Config{
			InsecureSkipVerify: true,
		}
	}

	p, err := c.GetPartition(config.Topic, cfg.ClientConfig)
	if err != nil {
		return nil, err
	}
	log.Info("get partitions", log.String("topic", config.Topic), log.Any("partition", p))
	// partitioned topic
	if p > 0 {
		list := make([]*consumerImpl, 0, p)
		originTopic := cfg.Topic
		for i := 0; i < p; i++ {
			cfg.Topic = fmt.Sprintf("%s-partition-%d", originTopic, i)
			mc := manage.NewManagedConsumer(c.pool, cfg)
			list = append(list, &consumerImpl{mcsm: mc, wg: config.Wg, closed: make(chan struct{})})
		}
		consumerList := &ConsumerList{
			list:             list,
			FlowPeriodSecond: DefaultFlowPeriodSecond,
			FlowPermit:       DefaultFlowPermit,
			closed:           make(chan struct{}),
		}
		return consumerList, nil
	}

	// single topic
	mc := manage.NewManagedConsumer(c.pool, cfg)
	return &consumerImpl{mcsm: mc, wg: config.Wg, closed: make(chan struct{})}, nil

}

func (c *client) reconnectProducer(config ProducerConfig) (*pub.Producer, error) {
	var err error
	var p *pub.Producer
	var i = 0
	for i < config.PublishRetry {
		p, err = c.newProducer(config)
		if err != nil {
			// failed and sleep
			time.Sleep(time.Duration(i+1) * time.Second)
		} else {
			// success and stop
			break
		}
		i++
	}

	if err != nil {
		log.Error("reconnect producer failed",
			log.Any("topic", config.Topic),
			log.Any("retry time", i+1),
			log.ErrorField(err),
		)
	} else {
		log.Info("reconnect producer success",
			log.Any("topic", config.Topic),
			log.Any("retry time", i+1),
		)
	}

	return p, err
}

func subscriptionName(topic string) string {
	return getTenant(topic) + "-sub"
}

func getTenant(topic string) string {
	topic = strings.TrimPrefix(topic, "persistent://")
	end := strings.Index(topic, "/")
	return topic[:end]
}

type clientPool struct {
	clients     []*client
	clientIndex int32
	poolSize    int32
}

func NewClient(cfg ClientConfig) Client {
	size := 20
	clients := make([]*client, 0, size)
	for i := 0; i < size; i++ {
		c := &client{
			pool:   manage.NewClientPool(),
			config: cfg,
		}
		clients = append(clients, c)
	}
	return &clientPool{
		clients:  clients,
		poolSize: int32(size),
	}
}

func (cp *clientPool) NewConsumer(config ConsumerConfig) (Consumer, error) {
	index := cp.clientIndex % cp.poolSize
	c, err := cp.clients[index].NewConsumer(config)
	atomic.AddInt32(&cp.clientIndex, 1)
	return c, err
}

func (cp *clientPool) NewProducer(config ProducerConfig) (Producer, error) {
	index := cp.clientIndex % cp.poolSize
	p, err := cp.clients[index].NewProducer(config)
	atomic.AddInt32(&cp.clientIndex, 1)
	return p, err
}

func (cp *clientPool) NewProducers(config ProducersConfig) (Producers, error) {
	index := cp.clientIndex % cp.poolSize
	ps, err := cp.clients[index].NewProducers(config)
	atomic.AddInt32(&cp.clientIndex, 1)
	return ps, err
}
