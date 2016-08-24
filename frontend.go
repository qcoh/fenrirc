package main

import (
	"fenrirc/irc"
	"fenrirc/mondrian"
)

// Frontend comprises the server and channel windows of an IRC connection.
type Frontend struct {
	serverWindow *mondrian.MessageBuffer
	channels     []*mondrian.MessageBuffer
}

// NewFrontend constructs a Frontend.
func NewFrontend(serverWindow *mondrian.MessageBuffer) *Frontend {
	return &Frontend{serverWindow: serverWindow, channels: []*mondrian.MessageBuffer{}}
}

// Server returns the server window.
func (f *Frontend) Server() irc.Appender {
	return f.serverWindow
}

// NewChannel returns a channel window.
func (f *Frontend) NewChannel(name string) irc.Appender {
	// name?
	f.channels = append(f.channels, NewMessageBuffer())
	return f.channels[len(f.channels)-1]
}
