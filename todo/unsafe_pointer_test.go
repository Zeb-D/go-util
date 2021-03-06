package todo

import (
	"fmt"
	"reflect"
	"testing"
	"unsafe"
)

// 相关指针操作
func IsInSlice(value interface{}, slice interface{}) bool {
	switch reflect.TypeOf(slice).Kind() {
	case reflect.Slice, reflect.Array:
		s := reflect.ValueOf(slice)
		for i := 0; i < s.Len(); i++ {
			if reflect.DeepEqual(value, s.Index(i).Interface()) {
				return true
			}
		}
	default:
		return false
	}
	return false
}

func IsInSliceNormal(value int, sli []int) bool {
	for _, v := range sli {
		if value == v {
			return true
		}
	}
	return false
}

func Test_IsInSlice(t *testing.T) {
	fmt.Println(IsInSlice(10, []int{1, 2, 4}))
	fmt.Println(IsInSlice(10, []int{1, 2, 4, 10}))
	fmt.Println(IsInSliceNormal(10, []int{1, 2, 4}))
	fmt.Println(IsInSliceNormal(10, []int{1, 2, 4, 10}))
}

func Benchmark_IsInSlice(b *testing.B) {
	fmt.Println(IsInSlice(10, []int{1, 2, 4}))
	fmt.Println(IsInSlice(10, []int{1, 2, 4, 10}))
}

// 原值与指针互换
func Test_Unsafe_Pointer(t *testing.T) {
	var value int64 = 5
	var p1 = &value
	var p2 = (*interface{})(unsafe.Pointer(p1))
	var p3 = (*int32)(unsafe.Pointer(p1))
	fmt.Println("p1:", p1)
	fmt.Println("*p1:", *p1)
	fmt.Println("&p1:", &p1)
	fmt.Println("p2:", p2)
	fmt.Println("*p3:", *p3)

	*p1 = 9212121231231212311
	fmt.Println("value:", value)
	fmt.Println("*p3:", *p3) // over flow int32
	fmt.Println("p2:", p2)
}

func Test_Memory_Address(t *testing.T) {

}
