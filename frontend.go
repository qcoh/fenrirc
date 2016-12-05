package main

import (
	"fenrirc/config"
	"fenrirc/irc"
	"io"
)

// Frontend comprises the server and channel windows of an IRC connection.
type Frontend struct {
	server   *Server
	channels []*Channel

	client io.Writer
	conf   *config.Server
}

// NewFrontend constructs a Frontend.
func NewFrontend(conf *config.Server, client io.Writer) *Frontend {
	return &Frontend{
		conf:     conf,
		client:   client,
		server:   NewServer(conf, client),
		channels: []*Channel{},
	}

}

// Server returns the server window.
func (f *Frontend) Server() irc.Appender {
	return f.server
}

// NewChannel returns a channel window.
func (f *Frontend) NewChannel(name string) irc.Channel {
	f.channels = append(f.channels, NewChannel(f.server, f.client, name))
	return f.channels[len(f.channels)-1]
}
