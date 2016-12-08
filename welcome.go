package main

import (
	"fenrirc/mondrian"
	"fenrirc/msg"
)

var (
	welcome *mondrian.MessageBuffer
)

func init() {
	welcome = mondrian.NewMessageBuffer()
	welcome.Append(msg.NewSimple("ᚠᛖᚾᚱᛁᚱᚲ"))
}

// NewMessageBuffer returns a *mondrian.MessageBuffer.
// When called the first time, it returns `firstMB`,
// which displays a welcome message.
func NewMessageBuffer() *mondrian.MessageBuffer {
	if welcome != nil {
		ret := welcome
		welcome = nil
		return ret
	}
	return mondrian.NewMessageBuffer()
}
