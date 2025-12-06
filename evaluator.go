package yap

import (
	"encoding/json"
	"strings"
)

type EvalContext struct {
	Json    any
	FuncMap map[string]Function
}

type Evaluator struct {
	expression Expr
}

func (e *Evaluator) Eval(data string) (interface{}, error) {
	var js map[string]any
	err := json.Unmarshal([]byte(data), &js)

	if err != nil {
		return nil, err
	}

	return e.expression.Eval(&EvalContext{
		Json:    js,
		FuncMap: BuiltinFunctions,
	})

}

func NewEvaluator(str string) (*Evaluator, error) {
	tokens, err := Tokenize(strings.NewReader(str))

	if err != nil {
		return nil, err
	}

	parser := NewParser(tokens)

	expr, err := parser.Parse()

	if err != nil {
		return nil, err
	}

	return &Evaluator{expression: expr}, nil
}

func Evaluate(str string, data string) (interface{}, error) {
	evaluator, err := NewEvaluator(str)

	if err != nil {
		return nil, err
	}

	result, err := evaluator.Eval(data)

	if err != nil {
		return nil, err
	}

	encoded, err := json.Marshal(result)

	if err != nil {
		return nil, err
	}

	return string(encoded), nil
}
