package yap

import (
	"bufio"
	"errors"
	"io"
	"math/big"
	"strings"
	"unicode"
)

// very simple charset
const (
	OpenParen   = '('
	CloseParen  = ')'
	Comma       = ','
	Dot         = '.'
	Equal       = '='
	Exclamation = '!'
	LessThan    = '<'
	GreaterThan = '>'
	Quote       = '"'
)

type TokenType int

const (
	Identifier TokenType = iota
	String
	Numeric
	Operator
	Punctuation
)

type Token struct {
	Type      TokenType
	Literal   string
	Numeric   *big.Float
	IsDecimal bool
}

type Tokenizer struct {
	reader *bufio.Reader
}

func (t *Tokenizer) ReadString() (string, error) {
	var sb strings.Builder

	for {
		c, _, err := t.reader.ReadRune()

		if err != nil {
			return "", err
		}

		if c == '"' {
			break
		}
		// handle escaping

		if c == '\\' {
			// peek at next string

			peek, _, err := t.reader.ReadRune()

			if err != nil {
				return "", err
			}

			switch peek {
			case 'n':
				{
					sb.WriteRune('\n') // write new line
				}
			case 'r': // carriage returns
				{
					sb.WriteRune('\r') // write carriage return
				}
			case '"':
				{
					sb.WriteRune('"') // write quote
				}
			case '\\':
				{
					sb.WriteRune('\\') // write slash
				}
			default:
				{
					return "", errors.New("unsupported escape sequence")
				}
			}
			continue
		}

		sb.WriteRune(c)
	}

	return sb.String(), nil
}

func (t *Tokenizer) readEquality(op rune) (*Token, error) {
	second, _, err := t.reader.ReadRune()

	if err == io.EOF {
		if op == '=' {
			return nil, errors.New("incomplete operator")
		}

		return &Token{
			Literal: string(op),
			Type:    Operator,
		}, nil
	}

	if err != nil {
		return nil, err
	}

	// weird syntax like '=>' or '=<' or '=!'
	if op == '=' && second != '=' {
		return nil, errors.New("unsupported equality operation")
	}

	if op != '=' && second != '=' {
		return &Token{
			Literal: string(op),
			Type:    Operator,
		}, nil
	}

	// validate the third rune is not '=', ('===' or '>==', '<==') etc.
	third, _, err := t.reader.ReadRune()

	if err != nil {
		return nil, err
	}

	if third == '=' {
		return nil, errors.New("unsupported equality operation")
	}

	return &Token{
		Literal: string(op) + "=",
		Type:    Operator,
	}, nil
}

func (t *Tokenizer) readNumeric() (*Token, error) {
	t.reader.UnreadRune()

	var literal strings.Builder
	var numeric strings.Builder

	isDecimal := false

	for {
		c, _, err := t.reader.ReadRune()

		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, err
		}

		if c >= '0' && c <= '9' {
			literal.WriteRune(c)
			numeric.WriteRune(c)
		} else if c == '_' || c == ',' {
			// should support commas so we can do a >= 1,000,000,000
			literal.WriteRune(c)
		} else if c == '.' {
			if isDecimal {
				// error: already a decimal
				return nil, errors.New("invalid numeric, already a decimal")
			}

			// mark as a decimal
			isDecimal = true

			literal.WriteRune(c)
			numeric.WriteRune(c)
		} else {
			// unread
			t.reader.UnreadRune()
			break // break out of loop
		}
	}

	f, ok := new(big.Float).SetString(numeric.String())

	if !ok {
		return nil, errors.New("failed to parse float")
	}

	return &Token{
		Type:      Numeric,
		Numeric:   f,
		Literal:   literal.String(),
		IsDecimal: isDecimal,
	}, nil

}

func (t *Tokenizer) ReadToken() (*Token, error) {
	r, _, err := t.reader.ReadRune()

	if err != nil {
		return nil, err
	}

	switch r {
	case OpenParen, CloseParen, Comma:
		{
			return &Token{
				Type:    Punctuation,
				Literal: string(r),
			}, nil
		}
	case Equal, Exclamation, LessThan, GreaterThan:
		{
			return t.readEquality(r)
		}
	case Quote:
		{
			// read a string up to next quote
			str, err := t.ReadString()
			if err != nil {
				return nil, err
			}
			return &Token{
				Type:    String,
				Literal: str,
			}, nil
		}

	default:
		{
			// discard whitespace
			if unicode.IsSpace(r) {
				return nil, nil
			}

			// try reading a numeric
			if r >= '0' && r <= '9' {
				// number
				return t.readNumeric()
			}

			// try reading an identifier

			return &Token{}, nil
		}
	}

}

func NewTokenizer(r io.Reader) *Tokenizer {
	return &Tokenizer{
		reader: bufio.NewReader(r),
	}
}
