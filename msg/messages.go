package msg

import (
	"fenrirc/mondrian"
	"strings"
	"time"
)

var (
	// NewSimple is the constructor for `msg.Simple`.
	NewSimple = newSimple
	// NewDefault is the constructor for `msg.Default`.
	NewDefault = newDefault
	// NewJoin is the constructor for `msg.Join`.
	NewJoin = newJoin
	// NewPrivate is the constructor for `msg.Private`.
	NewPrivate = newPrivate
	// NewReplyTopic is the constructor for `msg.ReplyTopic`.
	NewReplyTopic = newReplyTopic
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
	Raw string
	ToA time.Time
}

// newDefault constructs a default message.
func newDefault(raw string, toa time.Time) mondrian.Message {
	return Wrap(&Default{Raw: raw, ToA: toa})
}

// Draw draws message.
func (d *Default) Draw(r *mondrian.Region) {
	r.LPrintf("[%02d:%02d] ", d.ToA.Hour(), d.ToA.Minute())
	r.Xbase = r.Cx
	r.Printf("%s", d.Raw)
}

func nickFromPrefix(prefix string) string {
	if nickEnd := strings.Index(prefix, "!"); nickEnd != -1 {
		return prefix[0:nickEnd]
	}
	return prefix
}

// Join displays a join message.
type Join struct {
	Nick    string
	Channel string
	ToA     time.Time
}

func newJoin(prefix string, name string, toa time.Time) mondrian.Message {
	return Wrap(&Join{Nick: nickFromPrefix(prefix), Channel: name, ToA: toa})
}

// Draw draws the message.
func (j *Join) Draw(r *mondrian.Region) {
	r.LPrintf("[%02d:%02d] ", j.ToA.Hour(), j.ToA.Minute())
	r.Xbase = r.Cx
	r.Printf("%s has joined %s", j.Nick, j.Channel)
}

// Private displays a private message. (All messages to channels are private messages!)
type Private struct {
	Nick    string
	Content string
	ToA     time.Time
}

func newPrivate(prefix string, content string, toa time.Time) mondrian.Message {
	return Wrap(&Private{Nick: nickFromPrefix(prefix), Content: content, ToA: toa})
}

// Draw draws the message.
func (p *Private) Draw(r *mondrian.Region) {
	r.LPrintf("[%02d:%02d] ", p.ToA.Hour(), p.ToA.Minute())
	r.Xbase = r.Cx
	r.LPrintf("<%s> ", p.Nick)
	r.Printf("%s", p.Content)
}

// ReplyTopic displays a RPL_TOPIC message.
type ReplyTopic struct {
	Channel string
	Topic   string
	ToA     time.Time
}

func newReplyTopic(channel string, topic string, toa time.Time) mondrian.Message {
	return Wrap(&ReplyTopic{Channel: channel, Topic: topic, ToA: toa})
}

// Draw draws the message.
func (rt *ReplyTopic) Draw(r *mondrian.Region) {
	r.LPrintf("[%02d:%02d] ", rt.ToA.Hour(), rt.ToA.Minute())
	r.Xbase = r.Cx
	// TODO: set by?
	r.Printf("Topic for %s: %s", rt.Channel, rt.Topic)
}
