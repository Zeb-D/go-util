package common

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

//StringReverse support uft-8 string
func StringReverse(src string) string {
	bytes := []rune(src)
	for from, to := 0, len(bytes)-1; from < to; from, to = from+1, to-1 {
		bytes[from], bytes[to] = bytes[to], bytes[from]
	}
	return string(bytes)
}

func SubStringRange(src string, begin, end int) string {
	//you will be args check
	rs := []rune(src)
	return string(rs[begin:end])
}

func SubString(src string, begin int) string {
	//you will be args check
	rs := []rune(src)
	return string(rs[begin:])
}

func MapToString(a map[interface{}]interface{}) (string, error) {
	bs, err := json.Marshal(a)
	if err != nil {
		return "", nil
	}
	return string(bs), nil
}

func ToString(a interface{}) (string, error) {
	switch a.(type) {
	case uint:
		return strconv.FormatUint(uint64(a.(uint)), 10), nil
	case uint8:
		return strconv.FormatUint(uint64(a.(uint8)), 10), nil
	case uint16:
		return strconv.FormatUint(uint64(a.(uint16)), 10), nil
	case uint32:
		return strconv.FormatUint(uint64(a.(uint32)), 10), nil
	case uint64:
		return strconv.FormatUint(a.(uint64), 10), nil
	case int:
		return strconv.Itoa(a.(int)), nil
	case int8:
		return strconv.FormatInt(int64(a.(int8)), 10), nil
	case int16:
		return strconv.FormatInt(int64(a.(int16)), 10), nil
	case int32:
		return strconv.FormatInt(int64(a.(int32)), 10), nil
	case int64:
		return strconv.FormatInt(a.(int64), 10), nil
	case string:
		return a.(string), nil
	case float32:
		return strconv.FormatFloat(float64(a.(float32)), 'E', -1, 32), nil
	case float64:
		//	// 'b' (-ddddp±ddd，二进制指数)
		//// 'e' (-d.dddde±dd，十进制指数)
		//// 'E' (-d.ddddE±dd，十进制指数)
		//// 'f' (-ddd.dddd，没有指数)
		//// 'g' ('e':大指数，'f':其它情况)
		//// 'G' ('E':大指数，'f':其它情况)
		return strconv.FormatFloat(a.(float64), 'E', -1, 64), nil
	case bool:
		return strconv.FormatBool(a.(bool)), nil
	case map[interface{}]interface{}:
		return MapToString(a.(map[interface{}]interface{}))
	case fmt.Stringer:
		return a.(fmt.Stringer).String(), nil
	default:
		return "", errors.New(fmt.Sprintf("unsupported type params type:%v",
			reflect.TypeOf(a).Kind().String()))
	}
}

//	ToBytes 有4个字节的指针
func ToBytes(key string) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(key)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

//	ToInt8s 有4个字节的指针
func ToInt8s(key string) (is []int8, err error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err = enc.Encode(key)
	if err != nil {
		return nil, err
	}
	bs := buf.Bytes()
	is = make([]int8, len(bs))
	for _, value := range bs {
		is = append(is, int8(value))
	}
	return is, nil
}

func DecodeHexOrBase64(content string) ([]byte, error) {
	dat := []byte(content)
	isHex := true
	for _, v := range dat {
		if v >= 48 && v <= 57 || v >= 65 && v <= 70 || v >= 97 && v <= 102 {
			// isHex = true
		} else {
			isHex = false
			break
		}
	}
	if isHex {
		d, err := hex.DecodeString(content)
		if len(d) == 0 || err != nil {
			return base64.StdEncoding.DecodeString(content)
		}
		return d, err
	} else {
		return base64.StdEncoding.DecodeString(content)
	}
}
