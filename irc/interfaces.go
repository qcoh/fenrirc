package irc

import (
	"fenrirc/cmd"
	"fenrirc/mondrian"
)

// Appender is an interface for channel and server windows.
type Appender interface {
	Append(mondrian.Message)
}

// Channel is the interface for a channel.
type Channel interface {
	Appender
	SetTopic(string)
	HasNick(string) bool
	RemoveNick(string)
	InsertNick(string)
	SetNicks([]string)
	// more to follow
}

// Frontend is an interface for the widgets corresponding to a IRC connection.
type Frontend interface {
	Server(cmd.Handler) Appender
	NewChannel(string, cmd.Handler) Channel
	Sync(func())
	//Remove(Appender)
}
