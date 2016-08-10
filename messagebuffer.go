package main

import (
	"./mondrian"
)

var (
	firstMB *mondrian.MessageBuffer
)

func init() {
	firstMB = mondrian.NewMessageBuffer()
	// TODO: fill with welcome msg
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
