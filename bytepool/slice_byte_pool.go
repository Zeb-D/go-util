package bytepool

import (
	"sync/atomic"
)

var (
	minSize = 1024
	maxSize = 1073741824
)

type sliceBytePool struct {
	strategy        Strategy
	capToFreeBucket map[int]Bucket
	status          Status
}

func (bp *sliceBytePool) Get(size int) []byte {
	atomic.AddInt64(&bp.status.getCount, 1)

	ss := up2power(size)
	if ss < minSize {
		ss = minSize
	}
	bucket := bp.capToFreeBucket[ss]

	buf := bucket.Get(size)
	if buf == nil {
		buf = make([]byte, size, ss)
		return buf
	}

	atomic.AddInt64(&bp.status.hitCount, 1)
	atomic.AddInt64(&bp.status.sizeBytes, int64(-cap(buf)))
	return buf
}

func (bp *sliceBytePool) Put(buf []byte) {
	c := cap(buf)
	atomic.AddInt64(&bp.status.putCount, 1)
	atomic.AddInt64(&bp.status.sizeBytes, int64(c))

	size := down2power(c)
	if size < minSize {
		size = minSize
	}

	bucket := bp.capToFreeBucket[size]

	bucket.Put(buf)
}

func (bp *sliceBytePool) RetrieveStatus() Status {
	return Status{
		getCount:  atomic.LoadInt64(&bp.status.getCount),
		putCount:  atomic.LoadInt64(&bp.status.putCount),
		hitCount:  atomic.LoadInt64(&bp.status.hitCount),
		sizeBytes: atomic.LoadInt64(&bp.status.sizeBytes),
	}
}

// @return 范围为 [2, 4, 8, 16, ..., 1073741824]，如果大于等于1073741824，则直接返回n
func up2power(n int) int {
	if n >= maxSize {
		return n
	}

	var i uint32
	for ; n > (2 << i); i++ {
	}
	return 2 << i
}

// @return 范围为 [2, 4, 8, 16, ..., 1073741824]
func down2power(n int) int {
	if n < 2 {
		return 2
	} else if n >= maxSize {
		return maxSize
	}

	var i uint32
	for {
		nn := 2 << i
		if n > nn {
			i++
		} else if n == nn {
			return n
		} else if n < nn {
			return 2 << (i - 1)
		}
	}
}
