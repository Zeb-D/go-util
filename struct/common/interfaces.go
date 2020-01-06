package common

import (
	"errors"
	"reflect"
	"strings"
)

type ICompare interface {
	// Compare one interface{},
	// int return gt、eq、lt, if errors then failed
	Compare(b interface{}) (int, error)
}

func Compare(a interface{}, b interface{}) (int, error) {
	aType := reflect.TypeOf(a).String()
	bType := reflect.TypeOf(b).String()
	_, aok := a.(ICompare)
	_, bok := b.(ICompare)

	if !(aok && bok) && aType != bType {
		return -1, errors.New("cannot compare different type params")
	}

	switch a.(type) {
	case uint:
		if a.(uint) > b.(uint) {
			return 1, nil
		} else if a.(uint) < b.(uint) {
			return -1, nil
		} else {
			return 0, nil
		}
	case uint8:
		if a.(uint8) > b.(uint8) {
			return 1, nil
		} else if a.(uint8) < b.(uint8) {
			return -1, nil
		} else {
			return 0, nil
		}
	case uint16:
		if a.(uint16) > b.(uint16) {
			return 1, nil
		} else if a.(uint16) < b.(uint16) {
			return -1, nil
		} else {
			return 0, nil
		}
	case uint32:
		if a.(uint32) > b.(uint32) {
			return 1, nil
		} else if a.(uint32) < b.(uint32) {
			return -1, nil
		} else {
			return 0, nil
		}
	case uint64:
		if a.(uint64) > b.(uint64) {
			return 1, nil
		} else if a.(uint64) < b.(uint64) {
			return -1, nil
		} else {
			return 0, nil
		}
	case int:
		if a.(int) > b.(int) {
			return 1, nil
		} else if a.(int) < b.(int) {
			return -1, nil
		} else {
			return 0, nil
		}
	case int8:
		if a.(int8) > b.(int8) {
			return 1, nil
		} else if a.(int8) < b.(int8) {
			return -1, nil
		} else {
			return 0, nil
		}
	case int16:
		if a.(int16) > b.(int16) {
			return 1, nil
		} else if a.(int16) < b.(int16) {
			return -1, nil
		} else {
			return 0, nil
		}
	case int32:
		if a.(int32) > b.(int32) {
			return 1, nil
		} else if a.(int32) < b.(int32) {
			return -1, nil
		} else {
			return 0, nil
		}
	case int64:
		if a.(int64) > b.(int64) {
			return 1, nil
		} else if a.(int64) < b.(int64) {
			return -1, nil
		} else {
			return 0, nil
		}
	case string:
		return strings.Compare(a.(string), b.(string)), nil
	case float32:
		if a.(float32) > b.(float32) {
			return 1, nil
		} else if a.(float32) < b.(float32) {
			return -1, nil
		} else {
			return 0, nil
		}
	case float64:
		if a.(float64) > b.(float64) {
			return 1, nil
		} else if a.(float64) < b.(float64) {
			return -1, nil
		} else {
			return 0, nil
		}
	case bool:
		if (a.(bool) && b.(bool)) || (!a.(bool) && !b.(bool)) {
			return 0, nil
		} else if a.(bool) && !b.(bool) {
			return 1, nil
		} else {
			return -1, nil
		}
	case ICompare:
		return a.(ICompare).Compare(b)

	default:
		return -1, errors.New("unsupported type params")
	}
}
