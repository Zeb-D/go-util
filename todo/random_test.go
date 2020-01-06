package todo

import (
	"fmt"
	"github.com/Zeb-D/go-util/struct/list"
	"math/rand"
	"testing"
	"time"
)

func TestRandom(t *testing.T) {
	users := []string{"大鹏", "饮风", "胜军", "升勇", "文洲"}
	arrayList := list.NewArrayList(len(users))
	rand.Seed(time.Now().UnixNano())
	for true {
		for v := users[rand.Intn(len(users))]; !arrayList.Contains(v); {
			arrayList.AddFirst(v)
		}
		if len(users) == arrayList.Size() {
			break
		}
	}
	fmt.Println("人员随机列表->", arrayList)

	tasks := []string{"1", "2", "3", "4",
		"5", "6", "7", "8",
		"9", "10", "11", "12",
		"13"}
	tasklist := list.NewArrayList(len(tasks))
	for true {
		for v := tasks[rand.Intn(len(tasks))]; !tasklist.Contains(v); {
			tasklist.AddFirst(v)
		}
		if len(tasks) == tasklist.Size() {
			break
		}
	}

	fmt.Println("任务随机列表->", tasklist)

}
