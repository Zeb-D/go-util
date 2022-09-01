package main

import (
	"fmt"
	"github.com/Zeb-D/go-util/cache"
)

func main() {
	gc := cache.New(10).
		LFU().
		Build()
	gc.Set("key", "ok")

	v, err := gc.Get("key")
	if err != nil {
		panic(err)
	}
	fmt.Println("value:", v)
}
