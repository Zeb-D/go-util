package common

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"testing"
)

func TestCompressAndUncompress(t *testing.T) {
	CompressAndUncompress([]byte(`aaaaaaassss`))
}

func TestCompress(t *testing.T) {
	src := "爱我中华爱我中华爱中华中华中华我中华"
	dst := ZlibCompress([]byte(src))
	fmt.Printf("str:%s,len:%d,srcLen:%d \n", string(dst), len(dst), len(src))
	dst2 := ZlibUnCompress(dst)
	fmt.Println(string(dst2))
	assert.Equal(t, src, string(dst2))
}

func TestGzip(t *testing.T) {
	fName := "../testdata/MyFile.gz"
	var r *bufio.Reader
	fi, err := os.OpenFile(fName, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v, Can't open %s: error: %s\n", os.Args[0], fName,
			err)
		os.Exit(1)
	}
	//写
	gw := gzip.NewWriter(fi)
	gw.Write([]byte(`Hello bufIo`))
	gw.Flush()
	//读
	fz, err := gzip.NewReader(fi)
	if err != nil {
		r = bufio.NewReader(fi)
	} else {
		r = bufio.NewReader(fz)
	}

	for {
		line, err := r.ReadString('\n')
		if err != nil {
			fmt.Println("Done reading file")
			os.Exit(0)
		}
		fmt.Println(line)
	}
	os.Remove(fName)
}

func TestIC(t *testing.T) {
	log.SetFlags(log.Lmicroseconds)

	// 单元测试 case
	uidss := [][]uint32{
		{1, 2, 3, 18, 32, 100},
		{1, 2, 3, 18, 32},
		{1, 2, 3, 18},
		{1, 2, 3, 17},
		{1, 2, 3, 16},
		{1, 2, 3, 15, 16, 17, 18},
		{1, 2, 3, 15, 16, 17},
		{1, 2, 3, 15, 16},
		{1, 2, 3, 15},
		{1, 2, 3},
		{1, 2},
		{1},
	}

	var compressors []Compressor

	compressors = append(compressors, &OriginCompressor{})
	compressors = append(compressors, &OriginCompressor{ZlibExt: true})

	compressors = append(compressors, &LFCompressor{FB: 0})
	compressors = append(compressors, &LFCompressor{FB: 0, ZlibExt: true})
	compressors = append(compressors, &LFCompressor{FB: 2})
	compressors = append(compressors, &LFCompressor{FB: 4})
	compressors = append(compressors, &LFCompressor{FB: 4, ZlibExt: true})

	for _, c := range compressors {
		for _, uids := range uidss {
			log.Println("-----")
			log.Println("in uid len:", len(uids))

			Sort(uids)
			b := c.Compress(uids)
			log.Println("len(b):", len(b))

			uids2 := c.UnCompress(b)
			log.Println("out uid len:", len(uids2))

			// assert check
			if len(uids) != len(uids2) {
				panic(0)
			}
			for i := range uids {
				if uids[i] != uids2[i] {
					panic(0)
				}
			}
			log.Println("-----")
		}
	}
}

func marshalWrap(ids []uint32) (ret []byte) {
	log.Println("> sort.")
	Sort(ids)
	log.Println("< sort.")

	log.Println("> marshal.")
	//var oc OriginCompressor
	//ret = oc.Marshal(ids)

	var lfc LFCompressor
	lfc.FB = 4
	ret = lfc.Compress(ids)
	log.Println("< marshal.")

	log.Println("> zlib. len:", len(ret))
	ret = ZlibWrite(ret)
	log.Println("< zlib. len:", len(ret))
	return
}

func unmarshalWrap(b []byte) (ret []uint32) {
	b = ZlibRead(b)

	//var oc OriginCompressor
	//ret = oc.Unmarshal(b)

	var lfc LFCompressor
	lfc.FB = 4
	ret = lfc.UnCompress(b)
	return
}
