package todo

import (
	"fmt"
	"reflect"
	"testing"
	"unsafe"
)

//	实现原理是用 reflect.Type 得出来的信息来直接做反射，而不依赖于 reflect.ValueOf。
//	reflect.Value结构体 ，每次反射都需要malloc

func TestStructOffset(t *testing.T) {
	type testObj struct {
		f1 string
	}
	testObj_ := testObj{}
	testObjType_ := reflect.TypeOf(testObj_)
	fmt.Printf("%v,%v,%v,%v,%v \n", testObjType_.NumField(), testObjType_.Size(), testObjType_.Kind(), testObjType_.Align(), testObjType_.FieldAlign())
	for i := 0; i < testObjType_.NumField(); i++ {
		fmt.Printf("%v,%v,%v \n", testObjType_.Field(i).Name, testObjType_.Field(i).Offset, testObjType_.Field(i).Type)
	}
}

type f struct {
}

func (p *f) Run(a string) {

}

func TestReflectMethod(tt *testing.T) {
	p := f{}
	t := reflect.TypeOf(&p)
	fmt.Printf("f有%d个方法\n", t.NumMethod())

	m := t.Method(0)
	mt := m.Type
	fmt.Printf("%s方法有%d个参数\n", m.Name, mt.NumIn())
	for i := 0; i < mt.NumIn(); i++ {
		fmt.Printf("\t第%d个参数是%#v\n", i, mt.In(i).String())
	}
}

// SliceChunk 任意类型分片
// list： []T
// ret: [][]T
func SliceChunk(list interface{}, chunkSize int) (ret interface{}) {
	v := reflect.ValueOf(list)
	ty := v.Type() // []T

	// 先判断输入的是否是一个slice
	if ty.Kind() != reflect.Slice {
		fmt.Println("the parameter list must be an array or slice")
		return nil
	}

	// 获取输入slice的长度
	l := v.Len()

	// 计算分块之后的大小
	chunkCap := l/chunkSize + 1

	// 通过反射创建一个类型为[][]T的slice
	chunkSlice := reflect.MakeSlice(reflect.SliceOf(ty), 0, chunkCap)
	if l == 0 {
		return chunkSlice.Interface()
	}

	var start, end int
	for i := 0; i < chunkCap; i++ {
		end = chunkSize * (i + 1)
		if i+1 == chunkCap {
			end = l
		}
		// 将切片的append到chunk中
		chunkSlice = reflect.Append(chunkSlice, v.Slice(start, end))
		start = end
	}
	return chunkSlice.Interface()
}

func TestSliceChunk(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5, 6, 7, 8}
	s := SliceChunk(slice[1:], 3)
	fmt.Println(s)
	//看看这里面的类型
	st := reflect.TypeOf(slice)
	fmt.Printf("Elem:%v,Size:%v,Kind:%v,Align:%v,FieldAlign:%v \n", st.Elem(), st.Size(),
		st.Kind(), st.Align(), st.FieldAlign())
	st = reflect.TypeOf(s)
	fmt.Printf("Elem:%v,Size:%v,Kind:%v,Align:%v,FieldAlign:%v \n", st.Elem(), st.Size(),
		st.Kind(), st.Align(), st.FieldAlign())
	_, ok := s.([][]int)
	fmt.Println(ok)
}

func TestSlfReflect(t *testing.T) {
	i := "1111"
	eface := *(*emptyInterface)(unsafe.Pointer(&i))
	fmt.Printf("%v \n", eface)
}

type tflag uint8
type nameOff int32 // offset to a name
type typeOff int32 // offset to an *rtype
type textOff int32 // offset from top of text section
// rtype is the common implementation of most values.
// It is embedded in other struct types.
//
// rtype must be kept in sync with ../runtime/type.go:/^type._type.
type rtype struct {
	size       uintptr
	ptrdata    uintptr // number of bytes in the type that can contain pointers
	hash       uint32  // hash of type; avoids computation in hash tables
	tflag      tflag   // extra type information flags
	align      uint8   // alignment of variable with this type
	fieldAlign uint8   // alignment of struct field with this type
	kind       uint8   // enumeration for C
	// function for comparing objects of this type
	// (ptr to object A, ptr to object B) -> ==?
	equal     func(unsafe.Pointer, unsafe.Pointer) bool
	gcdata    *byte   // garbage collection data
	str       nameOff // string form
	ptrToThis typeOff // type for pointer to this type, may be zero
}

// emptyInterface is the header for an interface{} value.
type emptyInterface struct {
	typ  *rtype
	word unsafe.Pointer
}
