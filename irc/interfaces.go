package irc

import (
	"fenrirc/mondrian"
)

type Appender interface {
	Append(mondrian.Message)
}

type Frontend interface {
	Server() Appender
	//NewChannel(string) Appender
	//Remove(Appender)
	Logf(string, ...interface{})
}