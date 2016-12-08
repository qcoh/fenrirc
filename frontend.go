package main

import (
	"fenrirc/config"
	"fenrirc/irc"
)

// Frontend comprises the server and channel windows of an IRC connection.
type Frontend struct {
	server   *Server
	channels []*Channel
	conf     *config.Server
	syncfn   func(func())
}

// NewFrontend constructs a Frontend.
func NewFrontend(conf *config.Server, syncfn func(func())) *Frontend {
	return &Frontend{
		conf:     conf,
		server:   NewServer(conf),
		channels: []*Channel{},
		syncfn:   syncfn,
	}

}

// Server returns the server window.
func (f *Frontend) Server() irc.Appender {
	return f.server
}

// NewChannel returns a channel window.
func (f *Frontend) NewChannel(name string) irc.Channel {
	f.channels = append(f.channels, NewChannel(name))
	return f.channels[len(f.channels)-1]
}

// Sync synchronizes `ff` to the main goroutine.
func (f *Frontend) Sync(ff func()) {
	f.syncfn(ff)
}
