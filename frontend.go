package main

import (
	"fenrirc/config"
	"fenrirc/irc"
	"fenrirc/mondrian"
	"io"
)

type View interface {
	irc.Appender
	mondrian.InteractiveWidget
	Handler
}

type namedView struct {
	view View
	name string
}

// Frontend comprises the server and channel windows of an IRC connection.
type Frontend struct {
	views  []*namedView
	vindex int

	client io.Writer
	conf   *config.Server
}

// NewFrontend constructs a Frontend.
func NewFrontend(conf *config.Server, client io.Writer) *Frontend {
	return &Frontend{
		conf:   conf,
		client: client,
		views:  []*namedView{&namedView{view: NewServer(client), name: conf.Host}},
	}

}

// Server returns the server window.
func (f *Frontend) Server() irc.Appender {
	return f.views[0].view
}

// NewChannel returns a channel window.
func (f *Frontend) NewChannel(name string) irc.Channel {
	for _, v := range f.views {
		if v.name == name {
			return v.view.(irc.Channel) // ugh
		}
	}
	f.views = append(f.views, &namedView{view: NewChannel(f.views[0].view, f.client, name), name: name})
	return f.views[len(f.views)-1].view.(irc.Channel) // ugh^2
}

func (f *Frontend) next() View {
	if f.vindex+1 < len(f.views) {
		f.vindex++
		return f.views[f.vindex].view
	}
	return nil
}

func (f *Frontend) prev() View {
	if f.vindex > 0 {
		f.vindex--
		return f.views[f.vindex].view
	}
	return nil
}

func (f *Frontend) first() View {
	f.vindex = 0
	return f.views[f.vindex].view
}

func (f *Frontend) last() mondrian.InteractiveWidget {
	f.vindex = len(f.views) - 1
	return f.views[f.vindex].view
}
