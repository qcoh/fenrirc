package main

import (
	"errors"
	"strings"
)

var (
	errMalformed = errors.New("Malformed input")
)

// Command is a command from the user (prompt).
type Command struct {
	// TODO: Name this better
	Command string
	Params  []string
	Raw     string
}

func parse(raw string) (*Command, error) {
	if len(raw) == 0 {
		return nil, errMalformed
	}
	raw = strings.TrimRight(raw, "\r\n")

	ret := &Command{Raw: raw}

	if raw[0] != '/' {
		return ret, nil
	}

	s := strings.Split(raw[1:], " ")
	if len(s) <= 1 {
		ret.Command = s[0]
	} else {
		ret.Command, ret.Params = s[0], s[1:]
	}
	return ret, nil
}
