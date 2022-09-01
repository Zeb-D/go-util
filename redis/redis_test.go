package redis

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

var config = &Config{
	Host:     "127.0.0.1:6379",
	Password: "yd_redis",
	PoolSize: 1,
}

func TestDefaultRedisClient(t *testing.T) {
	DefaultRedisClient(config)
	var key = "go-util:config"
	v, err := json.Marshal(config)
	//ret, err := RedisClient.Set(key, v, 0).Result()
	//fmt.Printf("1ret:%s,err:%s,v:%s \n", ret, err, string(v))
	ret := Client.Get(key).Val()
	fmt.Printf("2ret:%s,err:%s \n", ret, err)
	assert.Equal(t, string(v), ret)
}

//test redis ZSet
func TestRedisSScan(t *testing.T) {
	DefaultRedisClient(config)
	var key = "go-util:all" //ZSet-key
	Client.Del(key)
	v, _ := json.Marshal(config)
	total := Client.SAdd(key, v, "2211", "222", "33322").Val()
	fmt.Printf("Zset Set key:%s,success:%d \n", key, total)
	scanCmd := Client.SScan(key, 0, "*", 2) //每次获取2个
	fmt.Printf("Zset Get key:%s,scanCmd:%s \n", key, scanCmd)
	iterator := scanCmd.Iterator()
	var valList = *new([]string)
	for iterator.Next() {
		valList = append(valList, strings.Trim(iterator.Val(), `"`))
	}
	fmt.Printf("val List:%s \n", valList)
}
