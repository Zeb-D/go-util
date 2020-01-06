package bytepool

func Init(strategy Strategy) {
	defaultPool = NewSliceBytePool(strategy)
}

func init() {
	Init(StrategyMultiSlicePoolBucket)
}

var defaultPool SliceBytePool

func Get(size int) []byte {
	return defaultPool.Get(size)
}

func Put(buf []byte) {
	defaultPool.Put(buf)
}

func RetrieveStatus() Status {
	return defaultPool.RetrieveStatus()
}

type SliceBytePool interface {
	// 功能类似于 make([]byte, <size>)
	Get(size int) []byte

	Put(buf []byte)

	RetrieveStatus() Status
}

type Status struct {
	getCount  int64
	putCount  int64
	hitCount  int64
	sizeBytes int64
}

type Strategy int

const (
	// 底层桶使用sync.Pool，内部的[]byte由sync.Pool决定何时释放
	StrategyMultiStdPoolBucket = iota + 1

	// 底层桶使用切片，内部的[]byte永远不会释放
	StrategyMultiSlicePoolBucket
)

type Bucket interface {
	// 桶内无满足条件的[]byte时，返回nil
	Get(size int) []byte

	Put(buf []byte)
}

func NewSliceBytePool(strategy Strategy) SliceBytePool {
	var capToFreeBucket map[int]Bucket

	switch strategy {
	case StrategyMultiStdPoolBucket:
		capToFreeBucket = make(map[int]Bucket)
		for i := minSize; i <= maxSize; i <<= 1 {
			capToFreeBucket[i] = NewStdPoolBucket()
		}
	case StrategyMultiSlicePoolBucket:
		capToFreeBucket = make(map[int]Bucket)
		for i := minSize; i <= maxSize; i <<= 1 {
			capToFreeBucket[i] = NewSliceBucket()
		}
	}

	return &sliceBytePool{
		strategy:        strategy,
		capToFreeBucket: capToFreeBucket,
	}
}
