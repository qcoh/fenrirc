package main

import (
	"fenrirc/cmd"
	"fenrirc/mondrian"
	"github.com/nsf/termbox-go"
	"sort"
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

func (t *topic) SetTopic(line string) {
	t.line = line
	mondrian.Draw(t)
}

type nicklist struct {
	*mondrian.Region
	mondrian.Visible
	nicks []string
}

func (nl *nicklist) Resize(r *mondrian.Region) {
	nl.Region = r
}

func (nl *nicklist) Draw() {
	nl.Clear()
	nl.Move(0, 0)

	for _, v := range nl.nicks {
		if nl.Cy > nl.Height {
			break
		}
		nl.LPrint(v)
		nl.Cx = 0
		nl.Cy++
	}
}

func (nl *nicklist) HasNick(n string) bool {
	i := sort.SearchStrings(nl.nicks, n)
	return i != len(nl.nicks) && nl.nicks[i] == n
}

func (nl *nicklist) RemoveNick(n string) {
	if i := sort.SearchStrings(nl.nicks, n); i != len(nl.nicks) && nl.nicks[i] == n {
		nl.nicks = append(nl.nicks[:i], nl.nicks[i+1:]...)
		mondrian.Draw(nl)
	}
}

func (nl *nicklist) InsertNick(n string) {
	i := sort.SearchStrings(nl.nicks, n)
	if i < len(nl.nicks) && nl.nicks[i] == n {
		return
	}
	nl.nicks = append(nl.nicks, "")
	copy(nl.nicks[i+1:], nl.nicks[i:])
	nl.nicks[i] = n
	mondrian.Draw(nl)
}

func (nl *nicklist) SetNicks(nicks []string) {
	sort.Strings(nicks)
	nl.nicks = nicks
	mondrian.Draw(nl)
}

// A Channel combines a messagebuffer with a topic line.
type Channel struct {
	cmd.Handler
	*mondrian.Box
	*topic
	*mondrian.MessageBuffer
	*nicklist

	name string
}

// NewChannel constructs a channel.
func NewChannel(name string, handler cmd.Handler) *Channel {
	ret := &Channel{
		Box:           mondrian.NewBox(),
		Handler:       handler,
		topic:         &topic{},
		MessageBuffer: NewMessageBuffer(),
		nicklist:      &nicklist{},
		name:          name,
	}
	ret.Children = []mondrian.Widget{ret.topic, ret.MessageBuffer, ret.nicklist}
	ret.ResizeFunc = func(r *mondrian.Region) []*mondrian.Region {
		return []*mondrian.Region{
			{X: r.X, Y: r.Y, Width: r.Width, Height: 1},
			{X: r.X, Y: r.Y + 1, Width: r.Width - 12, Height: r.Height - 1},
			{X: r.X + r.Width - 12, Y: r.Y + 1, Width: 12, Height: r.Height - 1},
		}
	}
	return ret
}

func (c *Channel) Draw() {
	c.Box.Draw()
}

func (c *Channel) Resize(r *mondrian.Region) {
	c.Box.Resize(r)
}

func (c *Channel) IsVisible() bool {
	return c.Box.IsVisible()
}

func (c *Channel) SetVisibility(v bool) {
	c.Box.SetVisibility(v)
}
