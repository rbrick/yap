package yap

type Expr struct{}

type BinOp struct {
	Left     Expr
	Operator string
	Right    Expr
}
type FuncCall struct {
	Name string
	Args []Expr
}

type Literal[T any] struct {
	Value T
}
