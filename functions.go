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

	"where": func(ctx *EvalContext, args []Expr) (interface{}, error) {
		if len(args) != 2 {
			return nil, fmt.Errorf("where function requires exactly 2 arguments")
		}
		arr, err := args[0].Eval(ctx)
		if err != nil {
			return nil, err
		}

		conditionExpr := args[1]

		if conditionExpr == nil {
			return nil, fmt.Errorf("where function requires a condition expression")
		}

		if conditionExpr.(*BinOp) == nil {
			return nil, fmt.Errorf("where function requires a binary operation as condition")
		}

		arrSlice, ok := arr.([]interface{})

		if !ok {
			return nil, fmt.Errorf("where function requires first argument to be an array")
		}

		matches := []any{}
		for _, item := range arrSlice {
			// inject special @ identifier for iterative elements
			jsonItem := map[string]any{
				"@": item,
			}

			conditionCtx := &EvalContext{
				Json:    jsonItem,
				FuncMap: ctx.FuncMap,
			}

			result, err := conditionExpr.Eval(conditionCtx)

			if err != nil {
				return nil, err
			}

			if result.(bool) {
				matches = append(matches, item)
			}
		}

		return matches, nil
	},
}
