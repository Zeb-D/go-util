package common

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"io"
)

//compress and uncompress

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
