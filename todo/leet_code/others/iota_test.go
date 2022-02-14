package others

import (
	"fmt"
	"testing"
)

const (
	a = iota
	b = iota
)
const (
	name = "menglu"
	c    = iota
	d    = iota
)
const (
	BASE = 1 << iota // iota为0，1左移0位 = 1
	e                //此时iota为1，1左移1位 = 2
	f                //此时iota为2，1左移2位 = 4
)

func TestIota(t *testing.T) {
	fmt.Println(a)
	fmt.Println(b)
	fmt.Println(c)
	fmt.Println(d)
	fmt.Println(BASE)
	fmt.Println(e)
	fmt.Println(f)
	fmt.Println(f == 4)

	str1 := []string{"a", "b", "c"}
	str2 := str1[1:]
	str2[1] = "new"
	fmt.Println(str1)
	str2 = append(str2, "z", "x", "y")
	fmt.Println(str1)
	fmt.Println(str2)
}
