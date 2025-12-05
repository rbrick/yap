package yap

import (
	"fmt"
	"math/big"
)

type Expr interface {
	Node() Expr

	Eval(ctx *EvalContext) (interface{}, error)
}

type Ident struct {
	Name string
}

func (i *Ident) Node() Expr {
	return i
}

func (i *Ident) Eval(ctx *EvalContext) (interface{}, error) {

	path, err := ParsePath(i.Name)

	if err != nil {
		return nil, err
	}

	return path.Resolve(ctx.Json)
}

type BinOp struct {
	Left     Expr
	Operator string
	Right    Expr
}

func (b *BinOp) Node() Expr {
	return b
}

func (b *BinOp) toBigFloat(i interface{}) (*big.Float, bool) {
	switch x := i.(type) {
	case *big.Float:
		return x, true
	case float64:
		return big.NewFloat(x), true
	case int64:
		return NewFloatFromInt(int(x)), true
	}

	return nil, false
}

func (b *BinOp) numericEval(left, right interface{}) (bool, error) {
	lNum, ok := b.toBigFloat(left)
	if !ok {
		return false, fmt.Errorf("left operand is not a number")
	}
	rNum, ok := b.toBigFloat(right)
	if !ok {
		return false, fmt.Errorf("right operand is not a number")
	}

	switch b.Operator {
	case "==":
		return lNum.Cmp(rNum) == 0, nil
	case "!=":
		return lNum.Cmp(rNum) != 0, nil
	case "<":
		return lNum.Cmp(rNum) == -1, nil
	case ">":
		return lNum.Cmp(rNum) == 1, nil
	case "<=":
		return lNum.Cmp(rNum) != 1, nil
	case ">=":
		return lNum.Cmp(rNum) != -1, nil
	default:
		return false, fmt.Errorf("unsupported operator: %s", b.Operator)
	}
}

func (b *BinOp) isNumeric(v interface{}) bool {
	_, ok := v.(*big.Float)
	return ok
}

func (b *BinOp) Eval(ctx *EvalContext) (interface{}, error) {
	left, err := b.Left.Eval(ctx)
	if err != nil {
		return nil, err
	}
	right, err := b.Right.Eval(ctx)
	if err != nil {
		return nil, err
	}

	switch b.Operator {
	case "==":
		if b.isNumeric(left) && b.isNumeric(right) {
			return b.numericEval(left, right)
		}
		return left == right, nil
	case "!=":
		if b.isNumeric(left) && b.isNumeric(right) {
			return b.numericEval(left, right)
		}
		return left != right, nil
	case "<", ">", "<=", ">=":
		return b.numericEval(left, right)

	default:
		return nil, fmt.Errorf("unsupported operator: %s", b.Operator)
	}

	return nil, fmt.Errorf("invalid operands for operator: %s", b.Operator)
}

type FuncCall struct {
	Name string
	Args []Expr
}

func (f *FuncCall) Node() Expr {
	return f
}

func (f *FuncCall) Eval(ctx *EvalContext) (interface{}, error) {
	if function, exists := ctx.FuncMap[f.Name]; exists {
		return function(ctx, f.Args)
	}
	return nil, fmt.Errorf("undefined function: %s", f.Name)
}

type Literal[T string | *big.Float] struct {
	Value T
}

func (l *Literal[T]) Node() Expr {
	return l
}

func (l *Literal[T]) Eval(ctx *EvalContext) (interface{}, error) {
	return l.Value, nil
}
