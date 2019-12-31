package todo

import (
	"fmt"
	"log"
	"testing"
	"time"

	"go.etcd.io/etcd/clientv3"
	"golang.org/x/net/context"
)

var (
	dialTimeout    = 5 * time.Second
	requestTimeout = 2 * time.Second
	endpoints      = []string{"http://localhost:2379"}
)

func TestEtcd(t *testing.T) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: dialTimeout,
	})

	if err != nil {
		log.Fatal(err)
	}

	defer cli.Close()

	key1, value1 := "testkey1", "value"

	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	_, err = cli.Put(ctx, key1, value1)
	cancel()
	if err != nil {
		log.Println("Put failed. ", err)
	} else {
		log.Printf("Put {%s:%s} succeed\n", key1, value1)
	}

	ctx, cancel = context.WithTimeout(context.Background(), requestTimeout)
	resp, err := cli.Get(ctx, key1)
	cancel()
	if err != nil {
		log.Println("Get failed. ", err)
		return
	}

	for _, kv := range resp.Kvs {
		log.Printf("Get {%s:%s} \n", kv.Key, kv.Value)
	}

	done := make(chan bool)

	go func() {
		wch := cli.Watch(context.Background(), key1)

		for item := range wch {
			for _, ev := range item.Events {
				log.Printf("Type:%s, key:%s, value:%s\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
			}
		}
	}()

	go func() {
		for cnt := 0; cnt < 11; cnt++ {
			value := fmt.Sprintf("%s%d", "value", cnt)
			_, err = cli.Put(context.Background(), key1, value)
			if err != nil {
				log.Println("Put failed. ", err)
			} else {
				log.Printf("Put {%s:%s} succeed\n", key1, value)
			}
		}
		done <- true
	}()

	dresp, err := cli.Delete(context.Background(), key1)
	if err != nil {
		log.Println("Delete failed. ", err)
		return
	}
	log.Println("->", dresp.Deleted)
	for _, kv := range dresp.PrevKvs {
		log.Printf("Deleted {%s:%s} \n", kv.Key, kv.Value)
	}

	<-done

	log.Println("Done!")
}
