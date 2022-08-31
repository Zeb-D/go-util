#### golang-lru

This provides the lru package which implements a fixed-size thread safe LRU cache. It is based on the cache in Groupcache.



#### implementation

lru cache base on sync.RWMutex + cache/simple/lru that non-thread safe,

so package cache's  `TwoQueueCache、ARCCache、Cache` is thread safe.



#### Example

```go
l, _ := NewLRU(128)
for i := 0; i < 256; i++ {
    l.Add(i, nil)
}
if l.Len() != 128 {
    panic(fmt.Sprintf("bad len: %v", l.Len()))
}
```