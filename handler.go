package main

import (
	"fenrirc/irc"
)

// A Handler is the interface implemented by everything reacting to user (prompt) input.
type Handler interface {
	Handle(*Command)
}

// ServerHandler sends user input to the server.
type ServerHandler struct {
	client *irc.Client
	next   Handler
}

// Handle reacts to a command.
func (sh *ServerHandler) Handle(cmd *Command) {
	switch cmd.Command {
	case "WHOIS":
	default:
		if sh.next != nil {
			sh.next.Handle(cmd)
		}
	}
}

// ChannelHandler sends user input relevant to a channel to the server.
type ChannelHandler struct {
	client *irc.Client
	name   string
	next   Handler
}

// Handle reacts to a command.
func (ch *ChannelHandler) Handle(cmd *Command) {
	switch cmd.Command {

	default:
		if ch.next != nil {
			ch.next.Handle(cmd)
		}
	}
}
