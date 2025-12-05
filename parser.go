package yap

import (
	"fmt"
	"math/big"
)

type Parser struct {
	tokens []*Token
	pos    int
}

func (p *Parser) currentToken() *Token {
	if p.pos >= len(p.tokens) {
		return nil
	}
	return p.tokens[p.pos]
}

func (p *Parser) advance() {
	if p.pos < len(p.tokens) {
		p.pos++
	}
}

func (p *Parser) peekToken() *Token {
	if p.pos+1 >= len(p.tokens) {
		return nil
	}
	return p.tokens[p.pos+1]
}

func (p *Parser) parseBinaryOp(left Expr) (Expr, error) {
	p.advance() // consume operator

	operand := p.currentToken()
	p.advance() // consume right

	right, err := p.parseExpression()

	if err != nil {
		return nil, err
	}

	binOp := &BinOp{
		Left:     left,
		Operator: operand.Literal,
		Right:    right,
	}

	return binOp, nil
}

func (p *Parser) parseFunctionCall(ident *Ident) (Expr, error) {
	funcCall := &FuncCall{
		Name: ident.Name,
		Args: []Expr{},
	}

	// consume identifier
	p.advance()
	// consume '('
	p.advance()

	for {
		token := p.currentToken() // current token == func2

		if token == nil {
			break
		}
		if token.Type == Punctuation && token.Literal == ")" {
			break
		}

		if token.Type == Punctuation && token.Literal == "," {
			p.advance() // consume ','
			continue
		}

		arg, err := p.parseExpression()
		if err != nil {
			return nil, err
		}

		funcCall.Args = append(funcCall.Args, arg)

		p.advance()

	}
	return funcCall, nil
}

func (p *Parser) parseOp(left Expr) (Expr, error) {

	nextToken := p.peekToken()

	if nextToken == nil {
		return left, nil
	}

	switch nextToken.Type {
	case BinaryOperator:
		return p.parseBinaryOp(left)
	}

	return left, nil
}

func (p *Parser) parseIdentifier() (Expr, error) {
	token := p.currentToken()
	if token == nil {
		return nil, fmt.Errorf("unexpected end of input")
	}

	ident := &Ident{Name: token.Literal}
	nextToken := p.peekToken()

	if nextToken != nil && nextToken.Type == Punctuation && nextToken.Literal == "(" {
		return p.parseFunctionCall(ident)
	}

	return ident, nil
}

func (p *Parser) parseLiteral() (Expr, error) {
	token := p.currentToken()
	if token == nil {
		return nil, fmt.Errorf("unexpected end of input")
	}

	switch token.Type {
	case String:
		return &Literal[string]{Value: token.Literal}, nil
	case Numeric:
		return &Literal[*big.Float]{Value: token.Numeric}, nil
	case Identifier:
		return p.parseIdentifier()
	default:
		return nil, fmt.Errorf("unexpected token: %s", token.Literal)
	}
}

func (p *Parser) parseExpression() (Expr, error) {

	token := p.currentToken()
	if token == nil {
		return nil, nil
	}

	switch token.Type {
	case Identifier:
		left, err := p.parseIdentifier()
		if err != nil {
			return nil, err
		}
		return p.parseOp(left)
	case String, Numeric:
		left, err := p.parseLiteral()
		if err != nil {
			return nil, err
		}
		return left, nil
	default:
		return nil, fmt.Errorf("unexpected token: %s", token.Literal)
	}
}

func (p *Parser) Parse() (Expr, error) {
	return p.parseExpression()
}

func NewParser(tokens []*Token) *Parser {
	return &Parser{tokens: tokens}
}
