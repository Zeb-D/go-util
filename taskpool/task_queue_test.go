package taskpool

import (
	"fmt"
	"github.com/Zeb-D/go-util/log"
	"testing"
	"time"
)

func TestChanTaskQueue(t *testing.T) {
	maxSize := 100
	q := NewChanTaskQueue(int64(maxSize))
	for i := 0; i < maxSize; i++ {
		go func(i int, q *chanTaskQueue) {
			ret, err := q.Add(&testR{})
			fmt.Println(i, " ->", ret, " ->", err)
			q.IsEmpty()
		}(i, q)
	}
	time.Sleep(1 * time.Second)
	go func(q *chanTaskQueue) {
		time.Sleep(800 * time.Millisecond)
		fmt.Println("a22->begin")
		r, err := q.Take()
		fmt.Println("a11->", r, err)
	}(q)
	ret, err := q.Offer(&testR{}, 800*time.Millisecond)
	fmt.Println("a12->", ret, err)
	time.Sleep(1 * time.Second)
	fmt.Println(len(q.workQueue))
}

type testR struct {
	name string
}

func (t *testR) Run() (bool, error) {
	time.Sleep(1 * time.Second)
	log.Info(" ->", log.Any("name", t.name))
	return true, nil
}
