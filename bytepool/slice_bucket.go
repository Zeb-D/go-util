package bytepool

import (
	"sync"
)

type SliceBucket struct {
	m    sync.Mutex
	core [][]byte
}

func NewSliceBucket() *SliceBucket {
	return new(SliceBucket)
}

func (b *SliceBucket) Get(size int) []byte {
	b.m.Lock()
	defer b.m.Unlock()
	if len(b.core) == 0 {
		return nil
	}
	buf := b.core[len(b.core)-1]
	b.core = b.core[:len(b.core)-1]
	return buf[0:size]
}

func (b *SliceBucket) Put(buf []byte) {
	b.m.Lock()
	defer b.m.Unlock()
	b.core = append(b.core, buf)
}
