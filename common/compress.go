package common

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
	"sort"
)

//compress and uncompress

// 具体使用见 LFCompressor 和 OriginCompressor
type Compressor interface {
	// 将整型切片压缩成二进制字节切片
	Compress(ids []uint32) (ret []byte)
	// 将二进制字节切片反序列化为整型切片
	// 反序列化后得到的整型切片，切片中整型的顺序和序列化之前保持不变
	UnCompress(b []byte) (ids []uint32)
}

func CompressAndUncompress(src []byte) {
	var in bytes.Buffer
	w := zlib.NewWriter(&in)
	w.Write(src)
	w.Close()
	fmt.Println(in.String())

	var out bytes.Buffer
	r, _ := zlib.NewReader(&in)
	io.Copy(&out, r)
	fmt.Println(out.String())
}

func ZlibCompress(src []byte) []byte {
	var in bytes.Buffer
	w := zlib.NewWriter(&in)
	w.Write(src)
	w.Close()
	return in.Bytes()
}

func ZlibUnCompress(compressSrc []byte) []byte {
	b := bytes.NewReader(compressSrc)
	var out bytes.Buffer
	r, _ := zlib.NewReader(b)
	io.Copy(&out, r)
	return out.Bytes()
}

func Sort(ids []uint32) {
	sort.Slice(ids, func(i, j int) bool {
		return ids[i] < ids[j]
	})
}

func resetBuf(b []byte) []byte {
	for i := 0; i < len(b); i++ {
		b[i] = 0
	}
	return b
}

func ZlibWrite(in []byte) []byte {
	var b bytes.Buffer
	w := zlib.NewWriter(&b)
	_, _ = w.Write(in)
	_ = w.Close()
	return b.Bytes()
}

func ZlibRead(in []byte) (ret []byte) {
	b := bytes.NewReader(in)
	r, _ := zlib.NewReader(b)
	ret, _ = ioutil.ReadAll(r)
	return
}

type LFCompressor struct {
	FB      uint32 // 用几个字节的 bit 表示跟随的数据
	ZlibExt bool   // 压缩之后，是否再用 zlib 进一步压缩

	oc OriginCompressor // FB 为0时，退化成使用 OriginCompressor
}

// 传入的整型切片必须是从小到大有序排列
func (lfc *LFCompressor) Compress(ids []uint32) (ret []byte) {
	if lfc.FB == 0 {
		ret = lfc.oc.Compress(ids)
		if lfc.ZlibExt {
			ret = ZlibWrite(ret)
		}
		return ret
	}

	lBuf := make([]byte, 4)
	fBuf := make([]byte, lfc.FB)

	maxDiff := 8 * lfc.FB

	var hasLeader bool
	var leader uint32
	var stage int
	for i := range ids {
		if !hasLeader {
			stage = 1
			leader = ids[i]
			hasLeader = true
			continue
		}

		diff := uint32(ids[i] - leader)

		if diff > maxDiff {
			binary.LittleEndian.PutUint32(lBuf, leader)
			ret = append(ret, lBuf...)
			ret = append(ret, fBuf...)

			resetBuf(fBuf)
			stage = 2
			leader = ids[i]
		} else {
			stage = 3
			fBuf[(diff-1)/8] = fBuf[(diff-1)/8] | (1 << byte((diff-1)%8))
		}
	}

	switch stage {
	case 1:
		binary.LittleEndian.PutUint32(lBuf, leader)
		ret = append(ret, lBuf...)
		dummy := make([]byte, lfc.FB)
		ret = append(ret, dummy...)
	case 2:
		binary.LittleEndian.PutUint32(lBuf, leader)
		ret = append(ret, lBuf...)
		dummy := make([]byte, lfc.FB)
		ret = append(ret, dummy...)
	case 3:
		binary.LittleEndian.PutUint32(lBuf, leader)
		ret = append(ret, lBuf...)
		ret = append(ret, fBuf...)
	}
	if lfc.ZlibExt {
		ret = ZlibWrite(ret)
	}
	return
}

func (lfc *LFCompressor) UnCompress(b []byte) (ids []uint32) {
	if lfc.ZlibExt {
		b = ZlibRead(b)
	}
	if lfc.FB == 0 {
		return lfc.oc.UnCompress(b)
	}

	isLeaderStage := true
	var item uint32
	var leader uint32
	var index uint32
	for {
		if isLeaderStage {
			leader = binary.LittleEndian.Uint32(b[index:])
			ids = append(ids, leader)
			isLeaderStage = false
			index += 4
		} else {
			for i := uint32(0); i < lfc.FB; i++ {
				for j := uint32(0); j < 8; j++ {
					if ((b[index+i] >> j) & 1) == 1 {
						item = leader + (i * 8) + j + 1
						ids = append(ids, item)
					}
				}
			}

			isLeaderStage = true
			index += lfc.FB
		}

		if int(index) == len(b) {
			break
		}
	}
	return
}

type OriginCompressor struct {
	ZlibExt bool // 压缩之后，是否再用 zlib 进一步压缩
}

// 并不强制要求整型切片有序
func (oc *OriginCompressor) Compress(ids []uint32) (ret []byte) {
	ret = make([]byte, len(ids)*4)
	for i, id := range ids {
		binary.LittleEndian.PutUint32(ret[i*4:], id)
	}
	if oc.ZlibExt {
		ret = ZlibWrite(ret)
	}
	return
}

func (oc *OriginCompressor) UnCompress(b []byte) (ids []uint32) {
	if oc.ZlibExt {
		b = ZlibRead(b)
	}
	n := len(b) / 4
	for i := 0; i < n; i++ {
		id := binary.LittleEndian.Uint32(b[i*4:])
		ids = append(ids, id)
	}
	return
}
