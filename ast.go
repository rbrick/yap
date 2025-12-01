package yap

import "math/big"

type Expr interface{}

type Ident struct {
	// The JSON path to follow
	Path []string
}

type BinOp struct {
	Left     Expr
	Operator string
	Right    Expr
}

type UnaryOp struct {
	Operator string
	Operand  Expr
}

type FuncCall struct {
	Name string
	Args []Expr
}

type Literal[T string | *big.Float] struct {
	Value T
}
