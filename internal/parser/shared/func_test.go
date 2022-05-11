package shared

import (
	"testing"
)

func TestCamel2Name(t *testing.T) {
	cases := []struct {
		Str      string
		Expected string
	}{
		{
			"HelloWorld",
			"hello_world",
		},
		{
			"Hello World",
			"hello world",
		},
		{
			"Hello我的World",
			"hello我的_world",
		},
		{
			"Hello 我的 World",
			"hello 我的 world",
		},
		{
			"Hello 我的World",
			"hello 我的_world",
		},
		{
			"Hello 我的world",
			"hello 我的world",
		},
		{
			"Hello 全角　World",
			"hello 全角　world",
		},
	}

	for i, one := range cases {
		got := camel2name(one.Str)
		if got != one.Expected {
			t.Errorf("#%d expected %q, got %q", i, one.Expected, got)
		}
	}
}
