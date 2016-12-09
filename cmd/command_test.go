package cmd

import (
	"testing"
)

func equal(a, b *Command) bool {
	if a.Method != b.Method {
		return false
	}
	if a.Raw != b.Raw {
		return false
	}
	if len(a.Params) != len(b.Params) {
		return false
	}
	for k := range a.Params {
		if a.Params[k] != b.Params[k] {
			return false
		}
	}
	return true
}

func TestParse(t *testing.T) {
	input := []string{
		"/CONNECT -Host foo -Port 123",
		"something, something",
	}
	expected := []*Command{
		{
			Method: "CONNECT",
			Params: []string{"-Host", "foo", "-Port", "123"},
			Raw:    input[0],
		},
		{
			Method: "",
			Params: []string{},
			Raw:    input[1],
		},
	}

	for k, v := range input {
		parsed, err := Parse(v)
		if err != nil {
			t.Error(err)
		}
		e := expected[k]
		if !equal(parsed, e) {
			t.Errorf("%+v != %+v\n", parsed, e)
		}
	}
}
