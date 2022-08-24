package common

import (
	"testing"
)

func TestAny2String(t *testing.T) {
	a := "123"
	var tests = []struct {
		a      any
		expect string
	}{
		{
			a:      int(1),
			expect: "1",
		}, {
			a:      int16(1),
			expect: "1",
		}, {
			a:      int32(1),
			expect: "1",
		}, {
			a:      uint(1),
			expect: "1",
		}, {
			a:      uint64(1),
			expect: "1",
		}, {
			a:      float32(1),
			expect: "1",
		}, {
			a:      1.333,
			expect: "1.333",
		}, {
			a:      "1.333",
			expect: "1.333",
		}, {
			a:      []interface{}{"123", "345"},
			expect: `["123","345"]`,
		}, {
			a: struct {
				name string `json:"name"`
			}{
				name: "yd",
			},
			expect: "{}",
		}, {
			a:      &a,
			expect: "123",
		},
	}
	for i, tt := range tests {
		actualV := Any2String(tt.a)
		if actualV != tt.expect {
			t.Log(i, tt.a, tt.expect, actualV)
			t.Fail()
		}
	}
}
