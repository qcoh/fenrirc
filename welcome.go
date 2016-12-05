package main

import (
	"fenrirc/mondrian"
	"fenrirc/msg"
)

var (
	welcome *Welcome
)

type Welcome struct {
	*mondrian.MessageBuffer
}

func (*Welcome) Status() string {
	return ""
}

func (*Welcome) Handle(*Command) {
}

func init() {
	welcome = &Welcome{mondrian.NewMessageBuffer()}
	welcome.Append(msg.NewSimple("ᚠᛖᚾᚱᛁᚱᚲ"))
}

// NewMessageBuffer returns a *mondrian.MessageBuffer.
// When called the first time, it returns `firstMB`,
// which displays a welcome message.
func NewMessageBuffer() *mondrian.MessageBuffer {
	if welcome != nil {
		ret := welcome.MessageBuffer
		welcome = nil
		return ret
	}
	return mondrian.NewMessageBuffer()
}
