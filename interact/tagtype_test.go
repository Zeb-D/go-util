package interact

import (
	"fmt"
	"reflect"
	"runtime"
	"testing"
)

type TagType struct {
	// tags
	field1 bool   "An important answer"
	field2 string "The name of the thing"
	field3 int    "How much there are"
}

func refTag(tt TagType, ix int) {
	ttType := reflect.TypeOf(tt)
	ixField := ttType.Field(ix)
	fmt.Printf("%v\n", ixField.Tag)
}

func TestTagType(t *testing.T) {
	tt := TagType{true, "Barak Obama", 1}
	for i := 0; i < 3; i++ {
		refTag(tt, i)
	}
}

type IntVector []int

func (v IntVector) Sum() (s int) {
	for _, x := range v {
		s += x
	}
	return
}

func TestRange(t *testing.T) {
	fmt.Println(IntVector{1, 2, 3, 4}.Sum())
}

type Base struct{}

func (Base) Magic() {
	fmt.Println("base magic")
}

func (self Base) MoreMagic() {
	self.Magic()
	self.Magic()
}

type Voodoo struct {
	Base
}

func (Voodoo) Magic() {
	fmt.Println("voodoo magic")
}

func TestMagic(t *testing.T) {
	v := new(Voodoo)
	v.Magic()
	v.MoreMagic()
	fmt.Printf("%v\n", runtime.NumCgoCall())
	runtime.SetFinalizer(v, func(obj *Voodoo) {
		fmt.Println("SetFinalizer")
	})
}

type Square struct {
	float32
}

func TestArea(t *testing.T) {
	s := &Square{1.22}
	fmt.Println(s.float32 * s.float32)
	m := make(map[string]interface{})
	m["aaa"] = "aaa"
	m["123"] = 123
	//v 是 varI 转换到类型 T 的值，ok 会是 true；
	//否则 v 是类型 T 的零值，ok 是 false，也没有运行时错误发生。
	if v, ok := m["aaa"].(int); ok || !ok { // checked type assertion
		fmt.Printf("int v:%v, ok:%v \n", v, ok)
	}
	for k, v := range m {
		SwithType(k)
		SwithType(v)
	}

}

func SwithType(v interface{}) {
	switch v.(type) {
	case int:
		fmt.Printf("int v:%s \n", v)
	case string:
		fmt.Printf("string v:%s \n", v)
	}
}
