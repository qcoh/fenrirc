package cmd

import (
	"errors"
	"strings"
)

var (
	errMalformed = errors.New("Malformed input")
)

// Command is a command from the user (prompt).
type Command struct {
	Method string
	Params []string
	Raw    string
}

// Parse parses a raw string into a command.
func Parse(raw string) (*Command, error) {
	if len(raw) == 0 {
		return nil, errMalformed
	}

	ret := &Command{Raw: raw}

	if raw[0] != '/' {
		return ret, nil
	}

	s := strings.Split(raw[1:], " ")
	if len(s) <= 1 {
		ret.Method = s[0]
	} else {
		ret.Method, ret.Params = s[0], s[1:]
	}
	return ret, nil
}
