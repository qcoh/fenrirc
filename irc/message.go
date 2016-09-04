package irc

import (
	"errors"
	"strings"
	"time"
)

type message struct {
	Prefix   string
	Command  string
	Params   []string
	Trailing string
	Raw      string
	ToA      time.Time
	From     string
}

var (
	errMalformed = errors.New("Malformed message")
)

// See: http://calebdelnay.com/blog/2010/11/parsing-the-irc-message-format-as-a-client
func parse(raw string) (*message, error) {
	raw = strings.TrimRight(raw, "\r\n")
	ret := &message{Raw: raw, ToA: time.Now()}

	prefixEnd := -1
	if strings.HasPrefix(raw, ":") {
		prefixEnd = strings.Index(raw, " ")
		if prefixEnd < 0 {
			return nil, errMalformed
		}
		ret.Prefix = raw[1:prefixEnd]
	}

	trailingStart := strings.Index(raw, " :")
	if trailingStart >= 0 {
		ret.Trailing = raw[trailingStart+2:]
	} else {
		trailingStart = len(raw)
	}

	cmdparams := strings.Split(raw[prefixEnd+1:trailingStart], " ")
	if len(cmdparams) > 0 {
		ret.Command, ret.Params = cmdparams[0], cmdparams[1:]
	} else {
		return nil, errMalformed
	}
	return ret, nil
}
