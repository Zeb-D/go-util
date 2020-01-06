package bytepool

import "go.uber.org/atomic"

type SharedSliceByte struct {
	Core  []byte
	pool  SliceBytePool
	count atomic.Uint32
}

type SharedSliceByteOption struct {
	pool SliceBytePool
}

var defaultSharedSliceByteOption = SharedSliceByteOption{
	pool: defaultPool,
}

type ModSharedSliceByteOption func(option *SharedSliceByteOption)

func WithPool(pool SliceBytePool) ModSharedSliceByteOption {
	return func(option *SharedSliceByteOption) {
		option.pool = pool
	}
}

func NewSharedSliceByte(size int, modOptions ...ModSharedSliceByteOption) *SharedSliceByte {
	option := defaultSharedSliceByteOption
	for _, fn := range modOptions {
		fn(&option)
	}

	var ssb SharedSliceByte
	ssb.Core = option.pool.Get(size)
	ssb.pool = option.pool
	ssb.count.Store(1)
	return &ssb
}

func WrapSharedSliceByte(b []byte, modOptions ...ModSharedSliceByteOption) *SharedSliceByte {
	option := SharedSliceByteOption{
		pool: defaultPool,
	}
	for _, fn := range modOptions {
		fn(&option)
	}

	var ssb SharedSliceByte
	ssb.Core = b
	ssb.pool = option.pool
	ssb.count.Store(1)
	return &ssb
}

func (ssb *SharedSliceByte) Ref() *SharedSliceByte {
	ssb.count.Inc()
	return ssb
}

func (ssb *SharedSliceByte) ReleaseIfNeeded() {
	if ssb.count.Dec() == 0 {
		ssb.pool.Put(ssb.Core)
	}
}
