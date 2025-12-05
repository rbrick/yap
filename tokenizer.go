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

	Multiplication = '*'
	Addition       = '+'
	Subtraction    = '-'
	Division       = '/'
)

type TokenType int

const (
	Identifier     TokenType = iota // 0
	String                          // 1
	Numeric                         // 2 - numeric
	UnaryOperator                   // 3
	BinaryOperator                  // 4
	Punctuation                     // 5
	WhiteSpace
)

func (tt TokenType) String() string {
	switch tt {
	case Identifier:
		return "Identifier"
	case String:
		return "String"
	case Numeric:
		return "Numeric"
	case UnaryOperator:
		return "UnaryOperator"
	case BinaryOperator:
		return "BinaryOperator"
	case Punctuation:
		return "Punctuation"
	case WhiteSpace:
		return "WhiteSpace"
	default:
		return "Unknown"
	}
}

type Token struct {
	Type      TokenType
	Literal   string
	Numeric   *big.Float
	IsDecimal bool
}

func (t *Token) String() string {
	return t.Literal
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
			Type:    BinaryOperator,
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
			Type:    BinaryOperator,
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
		Type:    BinaryOperator,
	}, nil
}

func (t *Tokenizer) readNumeric() (*Token, error) {
	t.reader.UnreadRune()

	var literal strings.Builder
	var numeric strings.Builder

	lastSeparator := false
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

			lastSeparator = false
		} else if c == '_' || c == ',' {
			if lastSeparator {
				// error: two separators in a row
				return nil, errors.New("invalid numeric, two separators in a row")
			}

			lastSeparator = true
			// should support commas so we can do a >= 1,000,000,000
			literal.WriteRune(c)
		} else if c == '.' {
			if isDecimal {
				// error: already a decimal
				return nil, errors.New("invalid numeric, already a decimal")
			}

			// mark as a decimal
			isDecimal = true
			lastSeparator = false

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

func (t *Tokenizer) readIdentifier(first rune) (*Token, error) {
	var literal strings.Builder
	literal.WriteRune(first)

	for {
		c, _, err := t.reader.ReadRune()

		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, err
		}

		if t.isIdentifierPart(c) {
			literal.WriteRune(c)
		} else {
			t.reader.UnreadRune()
			break
		}
	}

	return &Token{
		Type:    Identifier,
		Literal: literal.String(),
	}, nil
}

func (t *Tokenizer) isIdentifierStart(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' || r == '.' || r == '$' || r == '@'
}

func (t *Tokenizer) isIdentifierPart(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' || r == '.' || r == '$' || r == '[' || r == ']'
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
	case Multiplication, Addition, Subtraction, Division:
		{
			return &Token{
				Type:    BinaryOperator,
				Literal: string(r),
			}, nil
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
				return &Token{
					Type:    WhiteSpace,
					Literal: string(r),
				}, nil
			}

			// try reading a numeric
			if r >= '0' && r <= '9' {
				// number
				return t.readNumeric()
			}

			// try reading an identifier

			if t.isIdentifierStart(r) {
				return t.readIdentifier(r)
			}

			// unrecognized token
			return &Token{}, nil
		}
	}
}

func Tokenize(r io.Reader) ([]*Token, error) {
	tokenizer := NewTokenizer(r)
	var tokens []*Token

	for {
		token, err := tokenizer.ReadToken()

		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, err
		}

		// skip whitespace tokens
		if token.Type == WhiteSpace {
			continue
		}

		tokens = append(tokens, token)
	}

	return tokens, nil
}

func NewTokenizer(r io.Reader) *Tokenizer {
	return &Tokenizer{
		reader: bufio.NewReader(r),
	}
}
