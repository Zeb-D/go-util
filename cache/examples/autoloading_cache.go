package main

import (
	"fmt"
	"github.com/Zeb-D/go-util/cache"
)

func main() {
	gc := cache.New(10).
		LFU().
		LoaderFunc(func(key interface{}) (interface{}, error) {
			return fmt.Sprintf("%v-value", key), nil
		}).
		Build()

	v, err := gc.Get("key")
	if err != nil {
		panic(err)
	}
	fmt.Println(v)

	v2, err := gc.GetOrLoad("key2", func(key interface{}) (cache.Expirable, error) {
		return cache.Expirable{
			Value: fmt.Sprintf("%v-value", key),
		}, nil
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(v2)
}
