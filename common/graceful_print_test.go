package common

import "testing"

func TestPrint(t *testing.T) {
	bs := GracefulFprintf("hello", "world")
	println(string(bs))
}
