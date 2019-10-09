package redis

import (
	"fmt"
	"github.com/go-redis/redis"
	"time"
)

//you will init by DefaultRedisClient
var Client *redis.Client

//redis 配置文件
type Config struct {
	Host     string
	Password string
	DB       int
	PoolSize int
}

//生成全局默认的redis客户端
func DefaultRedisClient(config *Config) {
	Client = InitRedis(config)
}

//一个配置生成一个 客户端
func InitRedis(config *Config) *redis.Client {
	return config.initRedis()
}

//init redis
func (config Config) initRedis() *redis.Client {
	ops := &redis.Options{
		Addr:        config.Host,
		PoolSize:    config.PoolSize,  // Redis连接池大小
		MaxRetries:  3,                // 最大重试次数
		IdleTimeout: time.Second * 10, // 空闲链接超时时间
	}
	if len(config.Password) > 0 {
		ops.Password = config.Password
	}
	if config.DB != 0 {
		ops.DB = config.DB
	}
	redisClient := redis.NewClient(ops)
	if PingPong(redisClient) {
		fmt.Sprintf("redis init ping pong ok! \n")
	}
	go func() {
		for {
			time.Sleep(5 * time.Second)
			PingPong(redisClient)
		}
	}()
	return redisClient
}

//test for redis ping pong
func PingPong(redisClient *redis.Client) bool {
	//ping pong
	_, err := redisClient.Ping().Result()
	if err != nil {
		fmt.Errorf("redis init failed,ping has err:%s", err)
		return false
	}
	return true
}
