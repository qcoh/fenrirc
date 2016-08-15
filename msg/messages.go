package msg

import (
	"../mondrian"
)

var (
	NewSimple = newSimple
)

type message interface {
	Draw(*mondrian.Region)
}

// Cached stores the required height of a message depending on the width.
type Cached struct {
	message
	heightCache map[int]int
}

// Wrap returns a message with a height cache.
func Wrap(m message) *Cached {
	return &Cached{message: m, heightCache: make(map[int]int)}
}

// Height returns the height from the cache if it exists, otherwise computes, stores and returns it.
func (cm *Cached) Height(width int) int {
	if height, ok := cm.heightCache[width]; ok {
		return height
	}
	r := &mondrian.Region{X: 10000, Y: 10000, Width: width, Height: 10000}
	cm.Draw(r)
	cm.heightCache[width] = r.Cy + 1
	return r.Cy + 1
}

// Simple displays text.
type Simple struct {
	Text string
}

// NewSimple constructs a simple message.
func newSimple(text string) mondrian.Message {
	return Wrap(&Simple{text})
}

// Draw draws the message.
func (s *Simple) Draw(r *mondrian.Region) {
	r.LPrintf("%s", s.Text)
}
