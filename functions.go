package yap

import (
	"fmt"
)

type Function func(ctx *EvalContext, args []Expr) (interface{}, error)

var BuiltinFunctions = map[string]Function{
	"equals": func(ctx *EvalContext, args []Expr) (interface{}, error) {
		if len(args) != 2 {
			return nil, fmt.Errorf("equals function requires exactly 2 arguments")
		}
		left, err := args[0].Eval(ctx)
		if err != nil {
			return nil, err
		}
		right, err := args[1].Eval(ctx)
		if err != nil {
			return nil, err
		}
		return left == right, nil
	},

	"length": func(ctx *EvalContext, args []Expr) (interface{}, error) {
		if len(args) != 1 {
			return nil, fmt.Errorf("length function requires exactly 1 argument")
		}
		value, err := args[0].Eval(ctx)
		if err != nil {
			return nil, err
		}

		// we use big floats to represent numbers
		switch v := value.(type) {
		case string:
			return NewFloatFromInt(len(v)), nil
		case []interface{}:
			return NewFloatFromInt(len(v)), nil
		default:
			return nil, fmt.Errorf("length function not supported for type %T", value)
		}
	},
}
