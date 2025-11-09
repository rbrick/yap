package yap

import (
	"math/big"
	"strings"
	"testing"
)

func TestReadString(t *testing.T) {
	test := `"Hello, \"World\" 🔥🔥🔥🔥"` /// token.Literal == "Hello, \"World\" 🔥🔥🔥🔥"
	expect := "Hello, \"World\" 🔥🔥🔥🔥"
	tokenizer := NewTokenizer(strings.NewReader(test))

	str, err := tokenizer.ReadToken()

	if err != nil {
		t.FailNow()
	}

	if str.Literal != expect {
		t.Fail()
	}

	t.Log("read string:", str, "expected string:", expect)
}

func TestReadGreaterThan(t *testing.T) {
	test := `>`
	tokenizer := NewTokenizer(strings.NewReader(test))

	token, err := tokenizer.ReadToken()

	if err != nil {
		t.FailNow()
	}

	if token.Literal != ">" {
		t.Fail()
		t.Logf("failed to get greater than token: got %s, expected >", token.Literal)
	}
}

func TestReadNumeric(t *testing.T) {
	test := `10,000,000_000.314159 `
	expect := `10,000,000_000.314159`
	expectNumeric, _ := new(big.Float).SetString("10000000000.314159")
	tokenizer := NewTokenizer(strings.NewReader(test))

	token, err := tokenizer.ReadToken()

	if err != nil {
		t.FailNow()
	}

	if !token.IsDecimal {
		t.Fail()
		t.Logf("expected decimal, got false, result: %s", token.Literal)
	}

	if token.Literal != expect {
		t.Fail()
		t.Logf("got: %s, expected: %s", expect, token.Literal)
	}

	if token.Numeric.Cmp(expectNumeric) != 0 {
		t.Fail()
		t.Logf("expected: %s, got: %s", expectNumeric.String(), token.Numeric.String())
	}

}
