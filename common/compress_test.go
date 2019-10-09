package common

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"github.com/stretchr/testify/assert"
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
