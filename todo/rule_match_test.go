package todo

import (
	"github.com/nikunjy/rules/parser"
	"testing"
)

func TestEvaluate(t *testing.T) {
	type obj map[string]interface{}

	parser.Evaluate("x eq 1", obj{"x": 1})
	parser.Evaluate("x == 1", obj{"x": 1})
	parser.Evaluate("x lt 1", obj{"x": 1})
	parser.Evaluate("x < 1", obj{"x": 1})
	parser.Evaluate("x gt 1", obj{"x": 1})

	parser.Evaluate("x.a == 1 and x.b.c <= 2", obj{
		"x": obj{
			"a": 1,
			"b": obj{
				"c": 2,
			},
		},
	})

	parser.Evaluate("y == 4 and (x > 1)", obj{"x": 1})

	parser.Evaluate("y == 4 and (x IN [1,2,3])", obj{"x": 1})

	parser.Evaluate("y == 4 and (x eq 1.2.3)", obj{"x": "1.2.3"})
}

func TestOperations(t *testing.T) {

}
