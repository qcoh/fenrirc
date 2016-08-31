package main

import (
	"fenrirc/irc"
	"fenrirc/mondrian"
)

type namedView struct {
	view interface {
		irc.Appender
		mondrian.InteractiveWidget
	}
	name string
}

// Frontend comprises the server and channel windows of an IRC connection.
type Frontend struct {
	views  []*namedView
	vindex int
}

// NewFrontend constructs a Frontend.
func NewFrontend(host string) *Frontend {
	return &Frontend{views: []*namedView{{view: NewMessageBuffer(), name: host}}}
}

// Server returns the server window.
func (f *Frontend) Server() irc.Appender {
	return f.views[0].view
}

// NewChannel returns a channel window.
func (f *Frontend) NewChannel(name string) irc.Appender {
	for _, v := range f.views {
		if v.name == name {
			return v.view
		}
	}
	f.views = append(f.views, &namedView{view: NewMessageBuffer(), name: name})
	return f.views[len(f.views)-1].view
}

func (f *Frontend) next() mondrian.InteractiveWidget {
	if f.vindex+1 < len(f.views) {
		f.vindex++
		return f.views[f.vindex].view
	}
	return nil
}

func (f *Frontend) prev() mondrian.InteractiveWidget {
	if f.vindex > 0 {
		f.vindex--
		return f.views[f.vindex].view
	}
	return nil
}

func (f *Frontend) first() mondrian.InteractiveWidget {
	f.vindex = 0
	return f.views[f.vindex].view
}

func (f *Frontend) last() mondrian.InteractiveWidget {
	f.vindex = len(f.views) - 1
	return f.views[f.vindex].view
}
