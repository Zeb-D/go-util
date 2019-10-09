package interact

import (
	"container/list"
	"fmt"
	"log"
	"regexp"
	"runtime"
	"strconv"
	"testing"
	"time"
	"unicode/utf8"
	"unsafe"
)

func TestCase(t *testing.T) {
	fc := binOp(add)
	ret := fc(1, 3)
	fmt.Println(ret)
	a := callback(3, add)
	fmt.Println(a)
	fmt.Println(f())
}

func TestComplex64(t *testing.T) {
	var c1 complex64 = 5 + 10i
	fmt.Printf("c1:%v", c1)
}

func f() (ret int) {
	defer func() {
		ret++
	}()
	return 1
}

type binOp func(int, int) int

func add(x, y int) int {
	return x + y
}

func callback(a int, f func(int, int) int) int {
	return a + f(a, 2)
}

func fibonacci(n int) int {
	if n <= 1 {
		return 1
	}
	return fibonacci(n-1) + fibonacci(n-2)
}

func TestFibonacci(t *testing.T) {
	for i := 0; i <= 10; i++ {
		ret := fibonacci(i)
		fmt.Printf("i:%d,ret:%d \n", i, ret)
	}
}

func TestSlice(t *testing.T) {
	var arr1 [6]int
	var slice1 []int = arr1[2:5] // item at index 5 not included!

	// load the array with integers: 0,1,2,3,4,5
	for i := 0; i < len(arr1); i++ {
		arr1[i] = i
	}

	// print the slice
	for i := 0; i < len(slice1); i++ {
		fmt.Printf("Slice at %d is %d\n", i, slice1[i])
	}

	fmt.Printf("The length of arr1 is %d\n", len(arr1))
	fmt.Printf("The length of slice1 is %d\n", len(slice1))
	fmt.Printf("The capacity of slice1 is %d\n", cap(slice1))

	s := make([]byte, 5)
	s = s[1:3]
	fmt.Println(len(s))
	fmt.Println(cap(s))

}

func TestSliceCap(t *testing.T) {
	slice1 := make([]int, 0, 10)
	// load the slice, cap(slice1) is 10:
	for i := 0; i < cap(slice1); i++ {
		slice1 = slice1[0 : i+1]
		slice1[i] = i
		fmt.Printf("The length of slice is %d\n", len(slice1))
	}
}

func TestSliceCopy(t *testing.T) {
	sl_from := []int{1, 2, 3}
	sl_to := make([]int, 10)

	n := copy(sl_to, sl_from)
	fmt.Println(sl_to)
	fmt.Printf("Copied %d elements\n", n) // n == 3

	sl3 := []int{1, 2, 3}
	sl3 = append(sl3, 4, 5, 6)
	fmt.Println(sl3)
}

func TestStringSlice(t *testing.T) {
	s := "\u00ff\u754c"
	for i, c := range s {
		fmt.Printf("%d:%c ", i, c)
	}
	r := []rune(s)
	fmt.Println(r)
	fmt.Println(utf8.RuneCountInString(s))

	var b []byte
	var s1 string = "a,s,dfgg"
	b = append(b, s1...)
}

func TestList(t *testing.T) {
	l := list.List{}
	l.PushFront("aaa")
	for e := l.Front(); e != nil; e = e.Next() {
		fmt.Println(e.Value)
	}
	//多少个字节
	sizeof := unsafe.Sizeof(l)
	fmt.Println(sizeof)
	fmt.Println(unsafe.Sizeof("a"))
}

func TestRegexp(t *testing.T) {
	//目标字符串
	searchIn := "John: 2578.34 William: 4567.23 Steve: 5632.18"
	pat := "[0-9]+.[0-9]+" //正则

	f := func(s string) string {
		v, _ := strconv.ParseFloat(s, 32)
		return strconv.FormatFloat(v/2, 'f', 2, 32)
	}

	if ok, _ := regexp.Match(pat, []byte(searchIn)); ok {
		fmt.Println("Match Found!")
	}

	re, _ := regexp.Compile(pat)
	//将匹配到的部分替换为"##.#"
	str := re.ReplaceAllString(searchIn, "digits")
	fmt.Println(str)
	//参数为函数时
	str2 := re.ReplaceAllStringFunc(searchIn, f)
	fmt.Println(str2)
}

func TestFunc(t *testing.T) {

	where()
	// some code
	where()
}

var where func()

func init() {
	start := time.Now()
	where = func() {
		_, file, line, _ := runtime.Caller(1)
		log.Printf("%s:%d", file, line)
	}
	log.SetFlags(log.Llongfile)
	log.Print("当前耗时", time.Now().Sub(start))
}
