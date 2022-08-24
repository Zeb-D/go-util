package common

import (
	"encoding/json"
	"fmt"
)

func Any2String(a any) string {
	if a == nil {
		return ""
	}
	switch v := a.(type) {
	case uint, uint8, uint16, uint32, uint64:
		return fmt.Sprint(a)
	case int, int8, int16, int32, int64:
		return fmt.Sprint(a)
	case float32, float64:
		return fmt.Sprint(a)
	case string:
		return a.(string)
	case *string:
		return *v
	default:
		if bs, err := json.Marshal(a); err == nil {
			return string(bs)
		} else {
			return fmt.Sprint(a)
		}
	}
}
