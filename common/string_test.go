package common

import (
	"fmt"
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
