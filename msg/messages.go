package msg

import (
	"fenrirc/mondrian"
	"github.com/mattn/go-runewidth"
	"github.com/nsf/termbox-go"
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
	// NewNames is the constructor for `msg.Names`.
	NewNames = newNames
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

func drawTime(r *mondrian.Region, t time.Time) {
	r.AttrDefault()
	r.LPrintf("[%02d:%02d] ", t.Hour(), t.Minute())
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
func newDefault(m *Message) mondrian.Message {
	return Wrap(&Default{Raw: m.Raw, ToA: m.ToA})
}

// Draw draws message.
func (d *Default) Draw(r *mondrian.Region) {
	drawTime(r, d.ToA)
	r.Xbase = r.Cx
	r.Printf("%s", d.Raw)
}

func nickHostFromPrefix(prefix string) (string, string) {
	if nickEnd := strings.Index(prefix, "!"); nickEnd != -1 {
		return prefix[0:nickEnd], prefix[nickEnd+2:]
	}
	return prefix, ""
}

// Join displays a join message.
type Join struct {
	Nick    string
	Host    string
	Channel string
	ToA     time.Time
}

func newJoin(m *Message) mondrian.Message {
	n, h := nickHostFromPrefix(m.Prefix)
	var ch string
	if len(m.Params) > 0 {
		ch = m.Params[0]
	} else {
		ch = m.Trailing
	}

	return Wrap(&Join{Nick: n, Host: h, Channel: ch, ToA: m.ToA})
}

// Draw draws the message.
func (j *Join) Draw(r *mondrian.Region) {
	drawTime(r, j.ToA)
	r.Xbase = r.Cx
	r.Attr(termbox.ColorCyan|termbox.AttrBold, termbox.ColorDefault)
	r.Printf("%s", j.Nick)
	r.AttrDefault()
	r.LPrintf(" [")
	r.Attr(termbox.ColorCyan, termbox.ColorDefault)
	r.Printf("%s", j.Host)
	r.AttrDefault()
	r.Printf("] has joined")
	r.LPrintf(" ")
	r.Attr(termbox.ColorDefault|termbox.AttrBold, termbox.ColorDefault)
	r.Printf("%s", j.Channel)
	r.AttrDefault()
}

// Private displays a private message. (All messages to channels are private messages!)
type Private struct {
	Nick    string
	Content string
	ToA     time.Time
}

func newPrivate(m *Message) mondrian.Message {
	n, _ := nickHostFromPrefix(m.Prefix)
	return Wrap(&Private{Nick: n, Content: m.Trailing, ToA: m.ToA})
}

// Draw draws the message.
func (p *Private) Draw(r *mondrian.Region) {
	drawTime(r, p.ToA)
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

func newReplyTopic(m *Message) mondrian.Message {
	var ch string
	if len(m.Params) > 1 {
		ch = m.Params[1]
	}
	return Wrap(&ReplyTopic{Channel: ch, Topic: m.Trailing, ToA: m.ToA})
}

// Draw draws the message.
func (rt *ReplyTopic) Draw(r *mondrian.Region) {
	drawTime(r, rt.ToA)
	r.Xbase = r.Cx
	// TODO: set by?
	r.LPrintf("Topic for ")
	r.Attr(termbox.ColorCyan, termbox.ColorDefault)
	r.Printf("%s", rt.Channel)
	r.AttrDefault()
	r.Printf(": %s", rt.Topic)
}

// Names displays the users in the channel.
type Names struct {
	Names    []string
	MaxWidth int
	ToA      time.Time
}

func newNames(names []string, toa time.Time) mondrian.Message {
	maxwidth := 0
	for _, nick := range names {
		if w := runewidth.StringWidth(nick); w > maxwidth {
			maxwidth = w
		}
	}
	return Wrap(&Names{Names: names, MaxWidth: maxwidth + 2, ToA: toa})
}

// Draw draws the message.
func (n *Names) Draw(r *mondrian.Region) {
	r.Xbase = 8 // time
	ncol := (r.Width - r.Xbase) / n.MaxWidth
	if ncol == 0 {
		ncol = 1
	}
	nrow := len(n.Names)/ncol + 1
	if len(n.Names)%ncol == 0 {
		nrow--
	}

	drawRow := func(start int) {
		drawTime(r, n.ToA)
		for i := 0; start+i*nrow < len(n.Names); i++ {
			r.Cx = 8 + i*n.MaxWidth
			r.Attr(termbox.ColorDefault|termbox.AttrBold, termbox.ColorDefault)
			r.LPrintf("[")
			r.AttrDefault()
			r.LPrintf("%s", n.Names[start+i*nrow])
			r.Attr(termbox.ColorDefault|termbox.AttrBold, termbox.ColorDefault)
			r.Cx = 8 + (i+1)*n.MaxWidth - 1
			r.LPrintf("]")
			r.AttrDefault()
		}
	}

	for j := 0; j < nrow; j++ {
		r.Xbase = 0
		r.Cx = 0
		drawRow(j)
		r.Cy++
	}
	r.Cy--
}
