package main

import (
	"io"
)

// A Handler is the interface implemented by everything reacting to user (prompt) input.
type Handler interface {
	Handle(*Command)
}

// ServerHandler sends user input to the server.
type ServerHandler struct {
	client io.Writer
}

// Handle reacts to a command.
func (sh *ServerHandler) Handle(cmd *Command) {
	switch cmd.Command {
	case "WHOIS":
	}
}

// ChannelHandler sends user input relevant to a channel to the server.
type ChannelHandler struct {
	client io.Writer
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
