package etcd

import (
	"context"
	"github.com/Zeb-D/go-util/log"
	"strings"

	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"google.golang.org/grpc/resolver"
)

type Resolver struct {
	cli    *clientv3.Client
	cc     resolver.ClientConn
	scheme string
}

func NewResolver(EtcdCommonAddress string) resolver.Builder {
	return &Resolver{
		cli: NewEtcdClient(strings.Split(EtcdCommonAddress, ";")),
	}
}

func (r *Resolver) Scheme() string {
	return r.scheme
}

func (r *Resolver) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	go r.watch(r.scheme)
	return r, nil
}

// ResolveNow
func (r *Resolver) ResolveNow(rn resolver.ResolveNowOptions) {}

// Close
func (r *Resolver) Close() {}

func (r *Resolver) watch(prefix string) {
	addrDict := make(map[string]resolver.Address)

	update := func() {
		addrList := make([]resolver.Address, 0, len(addrDict))
		for _, v := range addrDict {
			addrList = append(addrList, v)
		}
		r.cc.NewAddress(addrList)
	}
	resp, err := r.cli.Get(context.Background(), prefix, clientv3.WithPrefix())
	if err == nil {
		for i, kv := range resp.Kvs {
			log.Info("resp", log.InfoField(i), log.InfoField(kv))
		}
	}

	update()

	rch := r.cli.Watch(context.Background(), prefix, clientv3.WithPrefix(), clientv3.WithPrevKV())
	for n := range rch {
		for _, ev := range n.Events {

			switch ev.Type {
			case mvccpb.PUT:
				log.Info("Watch Events", log.InfoField(ev))
			case mvccpb.DELETE:
				log.Info("Watch Events", log.InfoField(ev))
				delete(addrDict, string(ev.PrevKv.Key))
			}
		}
		update()
	}
}
