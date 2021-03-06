package etcd

import (
	"github.com/Zeb-D/go-util/log"
	"time"

	"github.com/coreos/etcd/clientv3"
)

func NewEtcdClient(endpoints []string) *clientv3.Client {
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: time.Second * 5,
	})
	if err != nil {
		log.Error("connect to etcd ", log.Any("endpoints", endpoints), log.ErrorField(err))
		panic(err)
	}
	log.Info("new etcd client success ", log.Any("endpoints", endpoints))
	return client
}
