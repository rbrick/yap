package yap

import (
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
