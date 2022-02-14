package todo

import (
	"encoding/json"
	jsoniter "github.com/json-iterator/go"
	"testing"
)

var jsonS = `{....300多个filed..}`

//	BenchmarkDefaultJSON-4   	 2419988	       467 ns/op
func BenchmarkDefaultJSON(b *testing.B) {
	for i := 0; i < b.N; i++ {
		param := make(map[string]interface{})
		_ = json.Unmarshal([]byte(jsonS), &param)
	}
}

//	BenchmarkIteratorJSON-4   	 1485985	       783 ns/op
func BenchmarkIteratorJSON(b *testing.B) {
	for i := 0; i < b.N; i++ {
		param := make(map[string]interface{})
		var json = jsoniter.ConfigCompatibleWithStandardLibrary
		_ = json.Unmarshal([]byte(jsonS), &param)
	}
}
