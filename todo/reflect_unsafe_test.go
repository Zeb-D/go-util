package todo

import (
	"fmt"
	"reflect"
	"testing"
)

// 定义结构体
type struct_ struct {
	// 可导出成员
	Name int
	// 不可导出成员
	age int
}

func TestType(t *testing.T) {
	// struct_的类型
	var type_struct = reflect.TypeOf(struct_{})
	fmt.Println("struct's kind:", type_struct.Kind(), ",name:", type_struct.Name())

	// 获取struct_的第一个成员类型
	var field_struct_Name, _ = type_struct.FieldByName("Name")
	fmt.Println("field_struct_name's kind:", field_struct_Name.Type.Kind(), ", name:", field_struct_Name.Type.Name())

	// slice
	var slices = []string{"hello", "world"}
	// 获取slice的类型
	var type_slice = reflect.TypeOf(slices)
	fmt.Println("slices's kind:", type_slice.Kind(), ", name:", type_slice.Name())
	// 获取slice中元素的类型
	var ele_slice = type_slice.Elem()
	fmt.Println("ele_slices's kind:", ele_slice.Kind(), ", name:", ele_slice.Name())

	// map
	var mmap = make(map[int]string)
	// 获取map的类型
	var type_mmap = reflect.TypeOf(mmap)
	fmt.Println("type_mmap's kind:", type_mmap.Kind(), ", name:", type_mmap.Name())
	// 获取map的key的类型
	var ele_key_mmap = type_mmap.Key()
	fmt.Println("ele_key_mmap's kind:", ele_key_mmap.Kind(), ", name:", ele_key_mmap.Name())
	// 获取map的value的类型
	var ele_value_mmap = type_mmap.Elem()
	fmt.Println("ele_value_mmap's kind:", ele_value_mmap.Kind(), ", name:", ele_value_mmap.Name())

	// *struct的类型
	var type_ptr = reflect.TypeOf(&struct_{})
	fmt.Println("type_ptr's kind:", type_ptr.Kind(), ", name:", type_ptr.Name())
	// *struct指向的结构体的类型
	var content_ptr = type_ptr.Elem()
	fmt.Println("content_ptr's kind:", content_ptr.Kind(), ", name:", content_ptr.Name())
}

func TestValue(t *testing.T) {
	// 通过结构体获取value值
	// 不可设置值，其成员也不能设置，哪怕是可导出成员
	var value_struct = reflect.ValueOf(struct_{})
	fmt.Println("canset of value_struct:", value_struct.CanSet())
	var ele_struct = value_struct.FieldByName("Name")
	fmt.Println("canset of ele_struct:", ele_struct.CanSet())

	// 通过*struct获取的value值可以进行设置值
	// 只有可导出成员可以设置值，不可导出成员不能设置新值
	var ptr_struct = &struct_{}
	var value_struct_ptr = reflect.ValueOf(ptr_struct)
	fmt.Println("canset of value_struct_ptr:", value_struct_ptr.CanSet())
	// 可导出成员可以设置值
	var field1_struct = value_struct_ptr.Elem().FieldByName("Name")
	fmt.Println("canset of field1_struct:", field1_struct.CanSet())
	field1_struct.SetInt(1000)
	// 不可导出成员不能设置值
	var field2_struct = value_struct_ptr.Elem().FieldByName("age")
	fmt.Println("canset of field2_struct:", field2_struct.CanSet())
	fmt.Println("the struct after set:", ptr_struct)

	// *slice 也可以进行值的设置
	var slice = []string{"hello", "world"}
	// *slice是不可设置的
	var value_slice_ptr = reflect.ValueOf(&slice)
	fmt.Println("canset of value_slice_ptr:", value_slice_ptr.CanSet())
	// *slice的元素是可以设置的
	var ele2_slice = value_slice_ptr.Elem().Index(1)
	fmt.Println("canset of ele2_slice:", ele2_slice.CanSet())
	ele2_slice.SetString("balabala")
	fmt.Println("slice after set:", slice)
	var slice_after_append = reflect.Append(value_slice_ptr.Elem(), reflect.ValueOf("world"))
	fmt.Println("slice after append:", slice_after_append)

	// map由于结构比较复杂，关于map中元素的修改不太可能
	var mmap = make(map[int]string)
	mmap[1] = "hello"
	mmap[2] = "world"
	var value_map_ptr = reflect.ValueOf(&mmap)
	fmt.Println("canset of value_map_ptr:", value_map_ptr.CanSet())

	var value_of_map = value_map_ptr.Elem().MapIndex(reflect.ValueOf(1))
	fmt.Println("canset of value_of_map:", value_of_map.CanSet())
}

//golang的基本类型的内存占用可以通过其底层数据类型得出。
//1. int 在64位机器中占用8个字节，在32位机器中占用4个字节
//2. int32占用4个字节
//3. int64占用8个字节
//4. float64占用8个字节
//5. float32占用4个字节
//6.指针或者uintptr在64位机器中占用8个字节，而在32位机器中占用4个字节
//7. string的底层类型是reflect.StringHeader，64位机器占用16个字节， 32位机器占用8个字节，定义如下。
//8. slice的底层类型是reflect.SliceHeader，64位机器占用24个字节，32位机器中占用12个字节，定义如下。
//9. bool占用一个字节
