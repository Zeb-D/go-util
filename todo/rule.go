package todo

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/caibirdme/yql"

	"github.com/nikunjy/rules/parser"
)

type Parser interface {
	Parse(rule, data []byte) (bool, error)
}

// json parser
type ruleParser struct{}

func NewRuleParser() Parser {
	return &ruleParser{}
}

func (p *ruleParser) Parse(rule, data []byte) (bool, error) {
	ev, err := parser.NewEvaluator(string(rule))
	if err != nil {
		log.Fatal(fmt.Errorf("Error making evaluator from the rule %v, %v", string(rule), err))
	}
	m := make(map[string]interface{})
	err = json.Unmarshal(data, &m)
	if err != nil {
		return false, err
	}
	ans, err := ev.Process(m)
	return ans, err
}

// sql parser
type yqlParser struct{}

func NewYqlParser() Parser {
	return &yqlParser{}
}

func (p *yqlParser) Parse(rule, data []byte) (bool, error) {
	m := make(map[string]interface{})
	err := json.Unmarshal(data, &m)
	if err != nil {
		return false, err
	}
	pass, err := yql.Match(string(rule), m)
	return pass, err
}
