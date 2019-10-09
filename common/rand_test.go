package common

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRandomCreateBytes(t *testing.T) {
	var a = 6
	a1 := RandomBytes(a)
	assert.Equal(t, len(a1), a)
	bs := []byte(`helloGoMyUtil`)
	_ = bs
	a2 := RandomBytes(a, []byte("helloGoMyUtil")...)
	assert.Equal(t, len(a1), len(a2))
	println("a1:", string(a1), "->a2:", string(a2))
}
