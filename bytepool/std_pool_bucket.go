package bytepool

import (
	"sync"
)

type StdPoolBucket struct {
	core *sync.Pool
}

func NewStdPoolBucket() *StdPoolBucket {
	return &StdPoolBucket{
		core: new(sync.Pool),
	}
}

func (b *StdPoolBucket) Get(size int) []byte {
	v := b.core.Get()
	if v == nil {
		return nil
	}
	vv := v.([]byte)
	return vv[0:size]
}

func (b *StdPoolBucket) Put(buf []byte) {
	b.core.Put(buf)
}
