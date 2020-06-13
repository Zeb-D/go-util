package common

import (
	"fmt"
	"github.com/Zeb-D/go-util/log"
	"github.com/pkg/errors"
	"reflect"
	"strings"
)

func DealGerResp(in interface{}, HumpToLine bool) (interface{}, error) {
	return dealResp(in, HumpToLine)
}

func dealResp(in interface{}, HumpToLine bool) (interface{}, error) {
	defer func() {
		if v := recover(); v != nil {
			log.Error("DealGerResp panic", log.Any("", v))
		}
	}()

	if in == nil {
		return in, nil
	}
	switch reflect.TypeOf(in).Kind() {
	case reflect.Map:
		if _, ok := in.(map[interface{}]interface{}); ok {
			m := MapIItoMapSI(in)
			if HumpToLine {
				m = Map2x_y(m)
			}
			return m, nil
		} else if inm, ok := in.(map[string]interface{}); ok {
			if HumpToLine {
				m := Map2x_y(in)
				return m, nil
			}
			return inm, nil
		} else {
			log.Error("can not find in data")
		}

	case reflect.Slice:
		value := reflect.ValueOf(in)
		newTemps := make([]interface{}, 0, value.Len())
		for i := 0; i < value.Len(); i++ {
			if value.Index(i).CanInterface() {
				newTemp, e := dealResp(value.Index(i).Interface(), HumpToLine)
				if e != nil {
					return nil, e
				}
				newTemps = append(newTemps, newTemp)
			} else {
				return nil, errors.New(fmt.Sprintf("unexpect err,value:%+v", value))
			}
		}
		return newTemps, nil
	default:
		return in, nil
	}
	return in, nil
}
func MapIItoMapSI(in interface{}) interface{} {
	inMap := make(map[interface{}]interface{})
	if v, ok := in.(map[interface{}]interface{}); !ok {
		return in
	} else {
		inMap = v
	}

	outMap := make(map[string]interface{}, len(inMap))
	for k, v := range inMap {
		if v == nil {
			continue
		}
		s := fmt.Sprint(k)
		if s == "class" {
			continue
		}

		vt := reflect.TypeOf(v)
		switch vt.Kind() {
		case reflect.Map:
			if _, ok := v.(map[interface{}]interface{}); ok {
				v = MapIItoMapSI(v)
			}
		case reflect.Slice:
			vl := reflect.ValueOf(v)
			os := make([]interface{}, 0, vl.Len())
			for i := 0; i < vl.Len(); i++ {
				if vl.Index(i).CanInterface() {
					osv := MapIItoMapSI(vl.Index(i).Interface())
					os = append(os, osv)
				}
			}
			v = os
		}
		outMap[s] = v

	}
	return outMap
}
func XY2_x_y(s string) string {
	data := make([]byte, 0, len(s)*2)
	j := false
	num := len(s)
	for i := 0; i < num; i++ {
		d := s[i]
		if i > 0 && d >= 'A' && d <= 'Z' && j {
			data = append(data, '_')
		}
		if d != '_' {
			j = true
		}
		data = append(data, d)
	}
	return strings.ToLower(string(data[:]))
}

/* map的key驼峰转下划线，会遍历所有的key */
func Map2x_y(in interface{}) interface{} {

	var m map[string]interface{}
	if v, ok := in.(map[string]interface{}); ok {
		m = v
	} else {
		return in
	}

	var out = make(map[string]interface{}, len(m))
	for k1, v1 := range m {
		x := XY2_x_y(k1)

		if v1 == nil {
			out[x] = v1
		} else if reflect.TypeOf(v1).Kind() == reflect.Struct {
			out[x] = Map2x_y(Struct2Map(v1))
		} else if reflect.TypeOf(v1).Kind() == reflect.Slice {
			value := reflect.ValueOf(v1)
			var newTemps = make([]interface{}, 0, value.Len())
			for i := 0; i < value.Len(); i++ {
				newTemp := Map2x_y(value.Index(i).Interface())
				newTemps = append(newTemps, newTemp)
			}
			out[x] = newTemps
		} else if reflect.TypeOf(v1).Kind() == reflect.Map {
			out[x] = Map2x_y(v1)
		} else {
			out[x] = v1
		}
	}
	return out
}
func Struct2MapAll(obj interface{}) interface{} {
	if obj == nil {
		return obj
	}
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)
	if t.Kind() == reflect.Struct {
		result := make(map[string]interface{})
		for i := 0; i < t.NumField(); i++ {
			if v.Field(i).Kind() == reflect.Struct {
				if v.Field(i).CanInterface() {
					setInMap(result, t.Field(i), Struct2MapAll(v.Field(i).Interface()))
				}
			} else if v.Field(i).Kind() == reflect.Slice || v.Field(i).Kind() == reflect.Map {
				if v.Field(i).CanInterface() {
					setInMap(result, t.Field(i), Struct2MapAll(v.Field(i).Interface()))
				}
			} else {
				if v.Field(i).CanInterface() {
					setInMap(result, t.Field(i), v.Field(i).Interface())
				}
			}
		}
		return result
	} else if t.Kind() == reflect.Slice {
		var newTemps = make([]interface{}, 0, v.Len())
		for i := 0; i < v.Len(); i++ {
			newTemp := Struct2MapAll(v.Index(i).Interface())
			newTemps = append(newTemps, newTemp)
		}
		return newTemps
	} else if t.Kind() == reflect.Map {
		var newTempMap = make(map[string]interface{}, v.Len())
		iter := v.MapRange()
		for iter.Next() {
			mapK := iter.Key().String()
			if !iter.Value().CanInterface() {
				continue
			}
			mapV := iter.Value().Interface()
			newTempMap[mapK] = Struct2MapAll(mapV)
		}
		return newTempMap
	} else {
		return obj
	}
}
func setInMap(m map[string]interface{}, structField reflect.StructField, value interface{}) (result map[string]interface{}) {
	result = m
	if tagName := structField.Tag.Get("m"); tagName == "" {
		result[headerAtoa(structField.Name)] = value
	} else {
		result[tagName] = value
	}
	return
}
func headerAtoa(a string) (b string) {
	b = strings.ToLower(a[:1]) + a[1:]
	return
}
func Struct2Map(obj interface{}) map[string]interface{} {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)

	var data = make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		data[t.Field(i).Name] = v.Field(i).Interface()
	}
	return data
}
