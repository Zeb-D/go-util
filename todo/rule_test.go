package todo

import (
	"fmt"
	"testing"
)

func TestRule(t *testing.T) {
	rule := "x < 2"
	data := `{"x":1}`
	p := NewRuleParser()
	pass, err := p.Parse([]byte(rule), []byte(data))
	fmt.Printf("pass:%v,err:%s \n", pass, err)
}
