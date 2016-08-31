package main

import (
	"fenrirc/mondrian"
	"github.com/nsf/termbox-go"
)

type topic struct {
	*mondrian.Region
	mondrian.Visible
	line string
}

func (t *topic) Resize(r *mondrian.Region) {
	t.Region = r
}

func (t *topic) Draw() {
	t.Clear()
	t.Move(0, 0)
	t.Attr(termbox.ColorDefault, termbox.ColorBlue)
	t.LPrintf("%s", t.line)
	t.AttrDefault()
}

// A Channel combines a messagebuffer with a topic line.
type Channel struct {
	*mondrian.Box

	t  *topic
	mb *mondrian.MessageBuffer
}

// NewChannel constructs a channel.
func NewChannel() *Channel {
	ret := &Channel{
		Box: mondrian.NewBox(),
		t:   &topic{},
		mb:  NewMessageBuffer(),
	}
	ret.Children = []mondrian.Widget{ret.t, ret.mb}
	ret.ResizeFunc = func(r *mondrian.Region) []*mondrian.Region {
		return []*mondrian.Region{
			{X: r.X, Y: r.Y, Width: r.Width, Height: 1},
			{X: r.X, Y: r.Y + 1, Width: r.Width, Height: r.Height - 1},
		}
	}
	return ret
}

// SetTopic sets the topic.
func (c *Channel) SetTopic(t string) {
	c.t.line = t
	mondrian.Draw(c)
}

// Append forwards msg to messagebuffer.
func (c *Channel) Append(msg mondrian.Message) {
	c.mb.Append(msg)
}

// HandleKey forwards ev to messagebuffer.
func (c *Channel) HandleKey(ev termbox.Event) {
	c.mb.HandleKey(ev)
}
