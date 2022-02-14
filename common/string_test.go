package common

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestStringReverse(t *testing.T) {
	start := time.Now().UnixNano()
	str := "我和我的祖国，我和我的家"
	fmt.Println(StringReverse(str))
	fmt.Println(time.Now().UnixNano() - start)
	fmt.Println(SubStringRange(str, 1, 5))
}

func TestStringIndex(t *testing.T) {
	s := "我和我的祖国，我和我的家"
	sub := "我"
	fmt.Println(strings.Index(s, sub))
	fmt.Println([]byte(s))
	fmt.Println([]byte("我"))
	sub = "我的"
	fmt.Println(strings.Index(s, sub))
}
