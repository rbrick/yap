package yap

import (
	"bufio"
	"errors"
	"io"
	"log"
	"strings"
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
	Type    TokenType
	Literal string
}

type Tokenizer struct {
	reader *bufio.Reader
}

func (t *Tokenizer) ReadString() (string, error) {
	var sb strings.Builder

	log.Println("hello world")

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
			return &Token{}, nil
		}
	}

}

func NewTokenizer(r io.Reader) *Tokenizer {
	return &Tokenizer{
		reader: bufio.NewReader(r),
	}
}
