package todo

import (
	"fmt"
	"testing"
)

func TestFX(t *testing.T) {
	fmt.Println(f1(1))
}

func f1(x int) (_, _ int) {
	_, _ = x, x
	return
}
