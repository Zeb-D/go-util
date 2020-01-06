package todo

import (
	"fmt"
	"reflect"
	"testing"
)

func TestReflect(t *testing.T) {
	i := 1
	v := reflect.ValueOf(&i)
	v.Elem().SetInt(10)
	//v.Set(11) panic
	fmt.Println(i)
}
