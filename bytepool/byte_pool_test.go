package bytepool

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewSharedSliceByte(t *testing.T) {
	pool := NewSliceBytePool(StrategyMultiSlicePoolBucket)
	assert.NotEmpty(t, pool)
	ssb := NewSharedSliceByte(1000, WithPool(pool))
	assert.NotEmpty(t, ssb)
	ssb2 := ssb.Ref()
	assert.NotEmpty(t, ssb2)

	ssb.ReleaseIfNeeded()
	ssb2.ReleaseIfNeeded()
}

func TestWrapSharedSliceByte(t *testing.T) {
	b := make([]byte, 1000)
	pool := NewSliceBytePool(StrategyMultiSlicePoolBucket)
	assert.NotEmpty(t, pool)
	ssb := WrapSharedSliceByte(b, WithPool(pool))
	assert.NotEmpty(t, ssb)
}

func TestDefault(t *testing.T) {
	Init(StrategyMultiSlicePoolBucket)
	buf := Get(1000)
	assert.Equal(t, 1000, len(buf))
	Put(buf)
	status := RetrieveStatus()
	e := Status{
		getCount:  1,
		putCount:  1,
		hitCount:  0,
		sizeBytes: 1024,
	}
	assert.Equal(t, e, status)
}

func TestMultiSlicePool(t *testing.T) {
	p := NewSliceBytePool(StrategyMultiSlicePoolBucket)
	assert.NotEmpty(t, p)
	buf := p.Get(1)
	assert.Equal(t, 1, len(buf))
	buf = make([]byte, 1)
	p.Put(buf)
	buf = p.Get(1)
	assert.Equal(t, 1, len(buf))
}

func TestMultiStdPool(t *testing.T) {
	p := NewSliceBytePool(StrategyMultiStdPoolBucket)
	assert.NotEmpty(t, p)
	buf := p.Get(1000)
	assert.Equal(t, 1000, len(buf))
	p.Put(buf)
	buf = p.Get(1000)
	assert.Equal(t, 1000, len(buf))
	for i := 0; i < 1000; i++ {
		buf = p.Get(1000)
		p.Put(buf)
	}
}

func TestUp2power(t *testing.T) {
	assert.Equal(t, 2, up2power(0))
	assert.Equal(t, 2, up2power(1))
	assert.Equal(t, 2, up2power(2))
	assert.Equal(t, 4, up2power(3))
	assert.Equal(t, 4, up2power(4))
	assert.Equal(t, 8, up2power(5))
	assert.Equal(t, 8, up2power(6))
	assert.Equal(t, 8, up2power(7))
	assert.Equal(t, 8, up2power(8))
	assert.Equal(t, 16, up2power(9))
	assert.Equal(t, 1024, up2power(1023))
	assert.Equal(t, 1024, up2power(1024))
	assert.Equal(t, 2048, up2power(1025))
	assert.Equal(t, 1073741824, up2power(1073741824-1))
	assert.Equal(t, 1073741824, up2power(1073741824))
	assert.Equal(t, 1073741824+1, up2power(1073741824+1))
	assert.Equal(t, 2047483647-1, up2power(2047483647-1))
	assert.Equal(t, 2047483647, up2power(2047483647))
}

func TestDown2power(t *testing.T) {
	assert.Equal(t, 2, down2power(0))
	assert.Equal(t, 2, down2power(1))
	assert.Equal(t, 2, down2power(2))
	assert.Equal(t, 2, down2power(3))
	assert.Equal(t, 4, down2power(4))
	assert.Equal(t, 4, down2power(5))
	assert.Equal(t, 4, down2power(6))
	assert.Equal(t, 4, down2power(7))
	assert.Equal(t, 8, down2power(8))
	assert.Equal(t, 8, down2power(9))
	assert.Equal(t, 512, down2power(1023))
	assert.Equal(t, 1024, down2power(1024))
	assert.Equal(t, 1024, down2power(1025))
	assert.Equal(t, 1073741824>>1, down2power(1073741824-1))
	assert.Equal(t, 1073741824, down2power(1073741824))
	assert.Equal(t, 1073741824, down2power(1073741824+1))
	assert.Equal(t, 1073741824, down2power(2047483647-1))
	assert.Equal(t, 1073741824, down2power(2047483647))
}
