package main

import (
	"fenrirc/mondrian"
	"fmt"
	"github.com/nsf/termbox-go"
	"time"
)

// StatusProvider is the interface providing a string for the Status widget.
type StatusProvider interface {
	Status() string
}

// TimeStatusProvider provides the current time.
type TimeStatusProvider struct{}

// Status returns the current time in hours and minutes.
func (TimeStatusProvider) Status() string {
	t := time.Now()
	return fmt.Sprintf("[%02d:%02d]", t.Hour(), t.Minute())
}

// Status is a widget which displays status information.
type Status struct {
	*mondrian.Region
	mondrian.Visible

	Global  StatusProvider
	Current StatusProvider
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
	s.Move(0, 0)

	for _, v := range []StatusProvider{s.Global, s.Current} {
		if v != nil {
			s.LPrint(v.Status())
			s.Cx++
		}
	}
}

// Resize resizes the Status widget.
func (s *Status) Resize(r *mondrian.Region) {
	s.Region = r
}
