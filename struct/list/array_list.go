package list

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/Zeb-D/go-util/struct/common"
)

type ArrayList struct {
	data []interface{}
	size int
}

// 构造函数，传入数组的容量capacity构造Array
func NewArrayList(capacity int) *ArrayList {
	return &ArrayList{
		data: make([]interface{}, capacity),
	}
}

// 获取数组的容量
func (arrayList *ArrayList) Capacity() int {
	return len(arrayList.data)
}

// 获得数组中的元素个数
func (arrayList *ArrayList) Size() int {
	return arrayList.size
}

// 返回数组是否为空
func (arrayList *ArrayList) IsEmpty() bool {
	return arrayList == nil || arrayList.size == 0
}

// 在第 index 个位置插入一个新元素 e
func (arrayList *ArrayList) Add(index int, e interface{}) error {
	if index < 0 || index > arrayList.size {
		return errors.New("Add failed. Require index >= 0 and index <= size.")
	}

	if arrayList.size == len(arrayList.data) {
		arrayList.resize(2 * arrayList.size)
	}

	for i := arrayList.size - 1; i >= index; i-- {
		arrayList.data[i+1] = arrayList.data[i]
	}

	arrayList.data[index] = e
	arrayList.size++
	return nil
}

// 向所有元素后添加一个新元素
func (arrayList *ArrayList) AddLast(e interface{}) error {
	return arrayList.Add(arrayList.size, e)
}

// 向所有元素前添加一个新元素
func (arrayList *ArrayList) AddFirst(e interface{}) error {
	return arrayList.Add(0, e)
}

// 获取 index 索引位置的元素
func (arrayList *ArrayList) Get(index int) interface{} {
	if index < 0 || index >= arrayList.size {
		panic("Get failed. Index is illegal.")
	}
	return arrayList.data[index]
}

// 修改 index 索引位置的元素
func (arrayList *ArrayList) Set(index int, e interface{}) {
	if index < 0 || index >= arrayList.size {
		panic("Set failed. Index is illegal.")
	}
	arrayList.data[index] = e
}

// 查找数组中是否有元素 e
func (arrayList *ArrayList) Contains(e interface{}) bool {
	for i := 0; i < arrayList.size; i++ {
		if ret, err := common.Compare(arrayList.data[i], e); err == nil && ret == 0 {
			return true
		}
	}
	return false
}

// 查找数组中元素 e 所在的索引，不存在则返回 -1
func (arrayList *ArrayList) Find(e interface{}) int {
	for i := 0; i < arrayList.size; i++ {
		if ret, err := common.Compare(arrayList.data[i], e); err == nil && ret == 0 {
			return i
		}
	}
	return -1
}

// 查找数组中元素 e 所有的索引组成的切片，不存在则返回 -1
func (arrayList *ArrayList) FindAll(e interface{}) (indexes []int) {
	for i := 0; i < arrayList.size; i++ {
		if ret, err := common.Compare(arrayList.data[i], e); err == nil && ret == 0 {
			indexes = append(indexes, i)
		}
	}
	return
}

// 从数组中删除 index 位置的元素，返回删除的元素
func (arrayList *ArrayList) Remove(index int) interface{} {
	if index < 0 || index >= arrayList.size {
		panic("Remove failed,Index is illegal.")
	}

	e := arrayList.data[index]
	for i := index + 1; i < arrayList.size; i++ {
		arrayList.data[i-1] = arrayList.data[i]
	}
	arrayList.size--
	arrayList.data[arrayList.size] = nil //loitering object != memory leak

	// 考虑边界条件，避免长度为 1 时，resize 为 0
	if arrayList.size == len(arrayList.data)/4 && len(arrayList.data)/2 != 0 {
		arrayList.resize(len(arrayList.data) / 2)
	}
	return e
}

// 从数组中删除第一个元素，返回删除的元素
func (arrayList *ArrayList) RemoveFirst() interface{} {
	return arrayList.Remove(0)
}

// 从数组中删除最后一个元素，返回删除的元素
func (arrayList *ArrayList) RemoveLast() interface{} {
	return arrayList.Remove(arrayList.size - 1)
}

// 从数组中删除一个元素 e
func (arrayList *ArrayList) RemoveElement(e interface{}) bool {
	index := arrayList.Find(e)
	if index == -1 {
		return false
	}

	arrayList.Remove(index)
	return true
}

// 从数组中删除所有元素 e
func (arrayList *ArrayList) RemoveAllElement(e interface{}) bool {
	if arrayList.Find(e) == -1 {
		return false
	}

	for i := 0; i < arrayList.size; i++ {
		if ret, err := common.Compare(arrayList.data[i], e); err == nil && ret == 0 {
			arrayList.Remove(i)
		}
	}
	return true
}

// 为数组扩容
func (arrayList *ArrayList) resize(newCapacity int) {
	newData := make([]interface{}, newCapacity)
	for i := 0; i < arrayList.size; i++ {
		newData[i] = arrayList.data[i]
	}

	arrayList.data = newData
}

func (arrayList *ArrayList) Swap(i int, j int) error {
	if i < 0 || i >= arrayList.size || j < 0 || j >= arrayList.size {
		return errors.New("Index is illegal.")
	}
	arrayList.data[i], arrayList.data[j] = arrayList.data[j], arrayList.data[i]
	return nil
}

// 重写 ArrayList 的 string 方法
func (arrayList *ArrayList) String() string {
	var buffer bytes.Buffer

	buffer.WriteString(fmt.Sprintf("ArrayList: size = %d, capacity = %d\n", arrayList.size, len(arrayList.data)))
	buffer.WriteString("[")
	for i := 0; i < arrayList.size; i++ {
		// fmt.Sprint 将 interface{} 类型转换为字符串
		buffer.WriteString(fmt.Sprint(arrayList.data[i]))
		if i != (arrayList.size - 1) {
			buffer.WriteString(", ")
		}
	}
	buffer.WriteString("]")

	return buffer.String()
}
