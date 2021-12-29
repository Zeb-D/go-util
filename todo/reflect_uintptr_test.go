package todo

import (
	"fmt"
	"reflect"
	"testing"
	"unsafe"
)

//此处reflect结合uintptr来高性能反射set值，因为 reflect.ValueOf 慢的一匹

//先解决一个小问题。怎么利用 reflect.StructField 取得对象上的值？
//对应的代码在： go/feature_reflect_object.go at master · json-iterator/go · GitHub

//在 reflect.StructField 上有一个 Offset 的属性。利用这个可以计算出字段的指针值。
func TestReflectStruct(t *testing.T) {
	type TestObj struct {
		field1 string
	}
	struct_ := &TestObj{}
	field, _ := reflect.TypeOf(struct_).Elem().FieldByName("field1")
	field1Ptr := uintptr(unsafe.Pointer(struct_)) + field.Offset
	*((*string)(unsafe.Pointer(field1Ptr))) = "hello"
	fmt.Println(struct_)
}

//获取 interface{} 的指针
//如果对应的结构体是以 interface{} 传进来的。还需要从 interface{} 上取得结构体的指针
func TestReflectInterface(t *testing.T) {
	type TestObj struct {
		field1 string
	}
	struct_ := &TestObj{}
	structInter := (interface{})(struct_)
	// emptyInterface is the header for an interface{} value.
	type emptyInterface struct {
		typ  *struct{}
		word unsafe.Pointer
	}
	structPtr := (*emptyInterface)(unsafe.Pointer(&structInter)).word
	field, _ := reflect.TypeOf(structInter).Elem().FieldByName("field1")
	field1Ptr := uintptr(structPtr) + field.Offset
	*((*string)(unsafe.Pointer(field1Ptr))) = "hello"
	fmt.Println(struct_)
}

//对应的代码在：go/feature_reflect_array.go at master · json-iterator/go · GitHub
func TestReflectSlice(t *testing.T) {
	slice := []string{"hello", "world", "china"}
	type sliceHeader struct {
		Data unsafe.Pointer
		Len  int
		Cap  int
	}
	header := (*sliceHeader)(unsafe.Pointer(&slice))
	fmt.Println(header.Len)
	elementType := reflect.TypeOf(slice).Elem()
	fmt.Printf("Kind:%v,Size:%v \n", elementType.Kind(), elementType.Size())
	secondElementPtr := uintptr(header.Data) + elementType.Size()
	*((*string)(unsafe.Pointer(secondElementPtr))) = "!!!"
	fmt.Println(slice)
	firstElementPtr := uintptr(header.Data)
	*((*string)(unsafe.Pointer(firstElementPtr))) = "```"
	fmt.Println(slice)
	threeElementPtr := secondElementPtr + elementType.Size()
	*((*string)(unsafe.Pointer(threeElementPtr))) = "aaa"
	fmt.Println(slice)
}

//Map :对于 Map 类型来说，没有 reflect.ValueOf 之外的获取其内容的方式。
//所以还是只能老老实实地用golang自带的值反射api。
func TestReflectMap(t *testing.T) {
	// map由于结构比较复杂，关于map中元素的修改不太可能
	var mmap = make(map[int]string)
	mmap[1] = "hello"
	mmap[2] = "world"
	var value_map = reflect.ValueOf(&mmap)
	fmt.Println("canset of value_map:", value_map.CanSet())

	var value_of_map = value_map.Elem().MapIndex(reflect.ValueOf(1))
	fmt.Println("canset of value_of_map:", value_of_map.CanSet())
	fmt.Println("Interface of value_of_map:", value_of_map.Interface())
}
