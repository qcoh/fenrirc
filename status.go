package main

import (
	"fenrirc/mondrian"
	"github.com/nsf/termbox-go"
)

// Status is a widget which displays status information.
type Status struct {
	*mondrian.Region
	mondrian.Visible
}

// NewStatus returns a Status widget.
func NewStatus() *Status {
	return &Status{
		Region: &mondrian.Region{Width: 100, Height: 100},
	}
}

// Draw draws status widget.
func (s *Status) Draw() {
	s.Bg = termbox.ColorBlue
	s.Clear()
}

// Resize resizes the Status widget.
func (s *Status) Resize(r *mondrian.Region) {
	s.Region = r
}
