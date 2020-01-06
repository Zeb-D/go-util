package todo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"
)

func TestVal(t *testing.T) {
	var a = 1
	a++
	x, y := add(&a), add(&a)
	fmt.Printf("x:%d,y:%d", x, y)
}

func add(a *int) int {
	b := *a
	b = b + 1
	a = &b
	return b
}

func TestJson(t *testing.T) {
	msg := `{"name":"zhangsan", "age":18, "id":92233720368547758, "sid":122464,"t":189.12}`

	var info map[string]interface{}
	decoder := json.NewDecoder(bytes.NewBuffer([]byte(msg)))
	decoder.UseNumber()
	if err := decoder.Decode(&info); err == nil {
		fmt.Println(info)
	}

	var info2 map[string]interface{}
	json.Unmarshal([]byte(msg), &info2)
	fmt.Println(info2)
}
