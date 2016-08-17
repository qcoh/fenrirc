package msg

import (
	"fenrirc/mondrian"
	"time"
)

var (
	// NewSimple is the constructor for `msg.Simple`.
	NewSimple = newSimple
	// NewDefault is the constructor for `msg.Default`.
	NewDefault = newDefault
	// NewLog is the constructor for `msg.Log`.
	NewLog = newLog
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

// newSimple constructs a simple message.
func newSimple(text string) mondrian.Message {
	return Wrap(&Simple{text})
}

// Draw draws the message.
func (s *Simple) Draw(r *mondrian.Region) {
	r.LPrintf("%s", s.Text)
}

// Default displays time, origin (irc network) and the raw line.
type Default struct {
	From string
	Raw  string
	ToA  time.Time
}

// newDefault constructs a default message.
func newDefault(from, raw string, toa time.Time) mondrian.Message {
	return Wrap(&Default{From: from, Raw: raw, ToA: toa})
}

// Draw draws message.
func (d *Default) Draw(r *mondrian.Region) {
	r.LPrintf("[%02d:%02d] ", d.ToA.Hour(), d.ToA.Minute())
	r.Xbase = r.Cx
	r.Printf("- %s - %s", d.From, d.Raw)
}

// Log displays a warning/log.
type Log struct {
	Text string
	From string
	ToA  time.Time
}

func newLog(text, from string, toa time.Time) mondrian.Message {
	return Wrap(&Log{Text: text, From: from, ToA: toa})
}

// Draw draws the message.
func (l *Log) Draw(r *mondrian.Region) {
	r.LPrintf("[%02d:%02d] ", l.ToA.Hour(), l.ToA.Minute())
	r.Xbase = r.Cx
	r.Printf(" - %s - %s", l.From, l.Text)
}
