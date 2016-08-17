package irc

import (
	"fenrirc/mondrian"
)

// Appender is an interface for channel and server windows.
type Appender interface {
	Append(mondrian.Message)
}

// Frontend is an interface for the widgets corresponding to a IRC connection.
type Frontend interface {
	Server() Appender
	//NewChannel(string) Appender
	//Remove(Appender)
}
