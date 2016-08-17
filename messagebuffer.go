package main

import (
	"fenrirc/mondrian"
	"fenrirc/msg"
)

var (
	firstMB *mondrian.MessageBuffer
)

func init() {
	firstMB = mondrian.NewMessageBuffer()
	firstMB.Append(msg.Wrap(&msg.Simple{"ᚠᛖᚾᚱᛁᚱᚲ"}))
}

// NewMessageBuffer returns a *mondrian.MessageBuffer.
// When called the first time, it returns `firstMB`,
// which displays a welcome message.
func NewMessageBuffer() *mondrian.MessageBuffer {
	if firstMB != nil {
		ret := firstMB
		firstMB = nil
		return ret
	}
	return mondrian.NewMessageBuffer()
}
