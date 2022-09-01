package cache

import (
	"errors"
	"fmt"
	"runtime/debug"
	"sync"
	"time"
)

const (
	TypeSimple = "simple"
	TypeLru    = "lru"
	TypeLfu    = "lfu"
	TypeArc    = "arc"
)

var KeyNotFoundError = errors.New("key not found")
var NoExpiration *time.Duration = nil

type Cache interface {
	Set(key, value interface{}) error
	SetWithExpire(key, value interface{}, expiration time.Duration) error
	Get(key interface{}) (interface{}, error)
	GetOrLoad(key interface{}, loader LoaderExpireFunc) (interface{}, error)
	GetIfPresent(key interface{}) (interface{}, error)
	GetAll(checkExpired bool) map[interface{}]interface{}
	Remove(key interface{}) bool
	Purge()
	Keys(checkExpired bool) []interface{}
	Len(checkExpired bool) int
	Has(key interface{}) bool

	get(key interface{}, onLoad bool) (interface{}, error)
	statsAccessor
}

type Expirable struct {
	Value  interface{}   // Value
	Expire time.Duration // Specified expiration of this value; expire<=0 will use default
}

type (
	LoaderFunc       func(interface{}) (interface{}, error)
	LoaderExpireFunc func(interface{}) (Expirable, error)
	EvictedFunc      func(interface{}, interface{})
	PurgeVisitorFunc func(interface{}, interface{})
	AddedFunc        func(interface{}, interface{})
	DeserializeFunc  func(interface{}, interface{}) (interface{}, error)
	SerializeFunc    func(interface{}, interface{}) (interface{}, error)
)

type Builder struct {
	clock            Clock
	tp               string
	size             int
	loaderExpireFunc LoaderExpireFunc
	evictedFunc      EvictedFunc
	purgeVisitorFunc PurgeVisitorFunc
	addedFunc        AddedFunc
	expiration       *time.Duration
	deserializeFunc  DeserializeFunc
	serializeFunc    SerializeFunc
}

type BaseCache struct {
	clock            Clock
	size             int
	loaderExpireFunc LoaderExpireFunc
	evictedFunc      EvictedFunc
	purgeVisitorFunc PurgeVisitorFunc
	addedFunc        AddedFunc
	deserializeFunc  DeserializeFunc
	serializeFunc    SerializeFunc
	expiration       *time.Duration
	mu               sync.RWMutex
	group            Group
	*stats
}

func NewExpirable(value interface{}, expire time.Duration) Expirable {
	return Expirable{
		Value: value, Expire: expire,
	}
}

func NewDefault(value interface{}) Expirable {
	return NewExpirable(value, time.Duration(-1))
}

func New(size int) *Builder {
	return NewWithEvictType(size, TypeSimple)
}

func NewLRU(size int) *Builder {
	return NewWithEvictType(size, TypeLru)
}

func NewLFU(size int) *Builder {
	return NewWithEvictType(size, TypeLfu)
}

func NewARC(size int) *Builder {
	return NewWithEvictType(size, TypeArc)
}

func NewWithEvictType(size int, tp string) *Builder {
	return &Builder{
		clock: NewRealClock(),
		tp:    tp,
		size:  size,
	}
}

func (b *Builder) Clock(clock Clock) *Builder {
	b.clock = clock
	return b
}

// LoaderFunc Set a loader function.
// loaderFunc: create a new value with this function if cached value is expired.
func (b *Builder) LoaderFunc(loaderFunc LoaderFunc) *Builder {
	return b.LoaderExpireFunc(func(k interface{}) (Expirable, error) {
		v, err := loaderFunc(k)
		return NewDefault(v), err
	})
}

// LoaderExpireFunc Set a loader function with expiration.
// loaderExpireFunc: create a new value with this function if cached value is expired.
// If nil returned instead of time.Duration from loaderExpireFunc than value will never expire.
func (b *Builder) LoaderExpireFunc(loaderExpireFunc LoaderExpireFunc) *Builder {
	b.loaderExpireFunc = loaderExpireFunc
	return b
}

func (b *Builder) EvictType(tp string) *Builder {
	b.tp = tp
	return b
}

func (b *Builder) Simple() *Builder {
	return b.EvictType(TypeSimple)
}

func (b *Builder) LRU() *Builder {
	return b.EvictType(TypeLru)
}

func (b *Builder) LFU() *Builder {
	return b.EvictType(TypeLfu)
}

func (b *Builder) ARC() *Builder {
	return b.EvictType(TypeArc)
}

func (b *Builder) EvictedFunc(evictedFunc EvictedFunc) *Builder {
	b.evictedFunc = evictedFunc
	return b
}

func (b *Builder) PurgeVisitorFunc(purgeVisitorFunc PurgeVisitorFunc) *Builder {
	b.purgeVisitorFunc = purgeVisitorFunc
	return b
}

func (b *Builder) AddedFunc(addedFunc AddedFunc) *Builder {
	b.addedFunc = addedFunc
	return b
}

func (b *Builder) DeserializeFunc(deserializeFunc DeserializeFunc) *Builder {
	b.deserializeFunc = deserializeFunc
	return b
}

func (b *Builder) SerializeFunc(serializeFunc SerializeFunc) *Builder {
	b.serializeFunc = serializeFunc
	return b
}

func (b *Builder) Expiration(expiration time.Duration) *Builder {
	b.expiration = &expiration
	return b
}

func (b *Builder) Build() Cache {
	if b.size <= 0 && b.tp != TypeSimple {
		panic("cache: Cache size <= 0")
	}

	return b.build()
}

func (b *Builder) build() Cache {
	switch b.tp {
	case TypeSimple:
		return newSimpleCache(b)
	case TypeLru:
		return newLRUCache(b)
	case TypeLfu:
		return newLFUCache(b)
	case TypeArc:
		return newARC(b)
	default:
		panic("cache: Unknown type " + b.tp)
	}
}

func buildCache(c *BaseCache, cb *Builder) {
	c.clock = cb.clock
	c.size = cb.size
	c.loaderExpireFunc = cb.loaderExpireFunc
	c.expiration = cb.expiration
	c.addedFunc = cb.addedFunc
	c.deserializeFunc = cb.deserializeFunc
	c.serializeFunc = cb.serializeFunc
	c.evictedFunc = cb.evictedFunc
	c.purgeVisitorFunc = cb.purgeVisitorFunc
	c.stats = &stats{}
}

// load a new value using by specified key.
func (c *BaseCache) load(key interface{}, exLoader LoaderExpireFunc, callback func(Expirable, error) (interface{}, error), isWait bool) (interface{}, bool, error) {
	v, called, err := c.group.Do(key, func() (v interface{}, e error) {
		defer func() {
			if r := recover(); r != nil {
				e = fmt.Errorf("[*bytepowered/base-cache] loader panics: %v, stack: %s", r, string(debug.Stack()))
			}
		}()
		return callback(exLoader(key))
	}, isWait)
	if err != nil {
		return nil, called, err
	}
	return v, called, nil
}
