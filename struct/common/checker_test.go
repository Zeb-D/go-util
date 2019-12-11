package common

import (
	"fmt"
	"testing"
)

func TestIsValid(t *testing.T) {
	fmt.Println(IsValid("()[]{}"))
	fmt.Println(IsValid("([)]"))
}
