# Cache

Golang缓存库，基于[lru](./lru)做特性增强，支持设定过期缓存，LFU, LRU, ARC, Simple等缓存算法实现。
同时也参考了[GCache](https://github.com/bluele/gcache)部分源码

## 特性

* 支持设定缓存失效时间：LFU, LRU and ARC；

* 支持协程安全；

* 支持监听缓存事件：evict, purge, added；(可选特性)

* 支持全局Loader自动加载不存在的键值；(可选特性)

* 支持外部化指定Loader加载不存在的键值；（通过`GetOrLoad`）

## 安装

```
$ go get github.com/Zeb-D/go-util
```

## 示例

### 手动设定Key-Value键值对

```go
package main

import (
  "github.com/Zeb-D/go-util/cache"
  "fmt"
)

func main() {
  gc := cache.New(20).
    LRU().
    Build()
  gc.Set("key", "ok")
  value, err := gc.Get("key")
  if err != nil {
    panic(err)
  }
  fmt.Println("Get:", value)
}
```

```
Get: ok
```

### 手动设定Key-Value键值对，指定缓存失效时间

```go
package main

import (
  "github.com/Zeb-D/go-util/cache"
  "fmt"
  "time"
)

func main() {
  gc := cache.New(20).
    LRU().
    Build()
  gc.SetWithExpire("key", "ok", time.Second*10)
  value, _ := gc.Get("key")
  fmt.Println("Get:", value)

  // Wait for value to expire
  time.Sleep(time.Second*10)

  value, err = gc.Get("key")
  if err != nil {
    panic(err)
  }
  fmt.Println("Get:", value)
}
```

```
Get: ok
// 10 seconds later, new attempt:
panic: ErrKeyNotFound
```


### 全局Loader加载不存在的键值

```go
package main

import (
  "github.com/Zeb-D/go-util/cache"
  "fmt"
)

func main() {
  gc := cache.New(20).
    LRU().
    LoaderFunc(func(key interface{}) (interface{}, error) {
      return "ok", nil
    }).
    Build()
  value, err := gc.Get("key")
  if err != nil {
    panic(err)
  }
  fmt.Println("Get:", value)
}
```

```
Get: ok
```

### 外部化指定Loader加载不存在的键值

```go
package main

import (
  "github.com/Zeb-D/go-util/cache"
  "fmt"
  "time"
)

func main() {
  gc := cache.New(20).
    LRU().
    Build()
  value, err := gc.GetOrLoad("key",func(key interface{}) (cache.Expirable, error) {
    return cache.NewDefault("my-new-value"), nil
  })
  if err != nil {
    panic(err)
  }
  fmt.Println("Get:", value)
}
```

```
Get: by-loader
```

### 全局Loader加载不存在的键值，并指定缓存失效时间

```go
package main

import (
  "fmt"
  "time"

  "github.com/Zeb-D/go-util/cache"
)

func main() {
  var evictCounter, loaderCounter, purgeCounter int
  gc := cache.New(20).
    LRU().
    LoaderExpireFunc(func(key interface{}) (cache.Expirable, error) {
      loaderCounter++
      return cache.NewExpirable("ok", time.Second*5), nil
    }).
    EvictedFunc(func(key, value interface{}) {
      evictCounter++
      fmt.Println("evicted key:", key)
    }).
    PurgeVisitorFunc(func(key, value interface{}) {
      purgeCounter++
      fmt.Println("purged key:", key)
    }).
    Build()
  value, err := gc.Get("key")
  if err != nil {
    panic(err)
  }
  fmt.Println("Get:", value)
  time.Sleep(1 * time.Second)
  value, err = gc.Get("key")
  if err != nil {
    panic(err)
  }
  fmt.Println("Get:", value)
  gc.Purge()
  if loaderCounter != evictCounter+purgeCounter {
    panic("bad")
  }
}
```

```
Get: ok
evicted key: key
Get: ok
purged key: key
```


## 缓存算法

  * Least-Frequently Used (LFU)

  Discards the least frequently used items first.

  ```go
  func main() {
    // size: 10
    gc := cache.New(10).
      LFU().
      Build()
    gc.Set("key", "value")
  }
  ```

  * Least Recently Used (LRU)

  Discards the least recently used items first.

  ```go
  func main() {
    // size: 10
    gc := cache.New(10).
      LRU().
      Build()
    gc.Set("key", "value")
  }
  ```

  * Adaptive Replacement Cache (ARC)

  Constantly balances between LRU and LFU, to improve the combined result.

  detail: http://en.wikipedia.org/wiki/Adaptive_replacement_cache

  ```go
  func main() {
    // size: 10
    gc := cache.New(10).
      ARC().
      Build()
    gc.Set("key", "value")
  }
  ```

  * SimpleCache (Default)

  SimpleCache has no clear priority for evict cache. It depends on key-value map order.

  ```go
  func main() {
    // size: 10
    gc := cache.New(10).Build()
    gc.Set("key", "value")
    v, err := gc.Get("key")
    if err != nil {
      panic(err)
    }
  }
  ```

## Loading Cache

If specified `LoaderFunc`, values are automatically loaded by the cache, and are stored in the cache until either evicted or manually invalidated.

```go
func main() {
  gc := cache.New(10).
    LRU().
    LoaderFunc(func(key interface{}) (interface{}, error) {
      return "value", nil
    }).
    Build()
  v, _ := gc.Get("key")
  // output: "value"
  fmt.Println(v)
}
```

GCache coordinates cache fills such that only one load in one process of an entire replicated set of processes populates the cache, then multiplexes the loaded value to all callers.

## Expirable cache

```go
func main() {
  // LRU cache, size: 10, expiration: after a hour
  gc := cache.New(10).
    LRU().
    Expiration(time.Hour).
    Build()
}
```

## 事件处理

### Evicted handler

Event handler for evict the entry.

```go
func main() {
  gc := cache.New(2).
    EvictedFunc(func(key, value interface{}) {
      fmt.Println("evicted key:", key)
    }).
    Build()
  for i := 0; i < 3; i++ {
    gc.Set(i, i*i)
  }
}
```

```
evicted key: 0
```

### Added handler

Event handler for add the entry.

```go
func main() {
  gc := cache.New(2).
    AddedFunc(func(key, value interface{}) {
      fmt.Println("added key:", key)
    }).
    Build()
  for i := 0; i < 3; i++ {
    gc.Set(i, i*i)
  }
}
```

```
added key: 0
added key: 1
added key: 2
```

# Authors

**Yd 2022 - 2023**

* <http://github.com/Zeb-D>
* <1406721322@qq.com>

**Jun Kimura @2019**

* <http://github.com/bluele>
* <junkxdev@gmail.com>
