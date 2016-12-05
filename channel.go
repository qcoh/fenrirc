package main

import (
	"fenrirc/mondrian"
	"fmt"
	"github.com/nsf/termbox-go"
	"io"
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
	t.Attr(termbox.ColorDefault, termbox.ColorBlue)
	t.Clear()
	t.Move(0, 0)
	t.LPrintf("%s", t.line)
	t.AttrDefault()
}

// A Channel combines a messagebuffer with a topic line.
type Channel struct {
	*mondrian.Box

	t  *topic
	mb *mondrian.MessageBuffer

	serverHandler Handler
	client        io.Writer
	name          string
}

// NewChannel constructs a channel.
func NewChannel(serverHandler Handler, client io.Writer, name string) *Channel {
	ret := &Channel{
		Box:           mondrian.NewBox(),
		t:             &topic{},
		mb:            NewMessageBuffer(),
		serverHandler: serverHandler,
		client:        client,
		name:          name,
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

// Handle handles user (prompt) input.
func (c *Channel) Handle(cmd *Command) {
	switch cmd.Command {
	case "":
		fmt.Fprintf(c.client, "PRIVMSG %s :%s\r\n", c.name, cmd.Raw)
		// TODO: notify that this msg was sent
	default:
		if c.serverHandler != nil {
			c.serverHandler.Handle(cmd)
		}
	}
}

// Status provides some channel info.
func (c *Channel) Status() string {
	// TODO
	return ""
}
