package common

import (
	"bytes"
	"fmt"
	"reflect"
)

//优雅输出对象 todo

func GracefulFprintf(data ...interface{}) string {
	var buf = new(bytes.Buffer)
	for _, value := range data {
		var buf2 = new(bytes.Buffer)
		KeyValue(buf2, reflect.ValueOf(value))
		fmt.Fprintln(buf2)
		buf.Write(buf2.Bytes())
	}
	return buf.String()
}

func KeyValue(buf *bytes.Buffer, val reflect.Value) {
	switch val.Kind() {
	case reflect.Bool:
		fmt.Println("val", val)
		fmt.Fprint(buf, val.Bool())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		fmt.Fprint(buf, val.Int())
	case reflect.String:
		fmt.Fprint(buf, "\"", val.String(), "\"")
	}
}
