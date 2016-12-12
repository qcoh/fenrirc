package irc

import (
	"fenrirc/cmd"
	"sort"
)

type channel struct {
	Channel
	*server
	name      string
	nicks     []string
	nicksTemp []string
}

func (c *channel) Handle(command *cmd.Command) {
	switch command.Method {
	case "":
		c.Writef("PRIVMSG %s :%s\r\n", c.name, command.Raw)
	default:
		c.server.Handle(command)
	}
}

func (c *channel) hasNick(n string) bool {
	i := sort.SearchStrings(c.nicks, n)
	return i != len(c.nicks) && c.nicks[i] == n
}

func (c *channel) removeNick(n string) {
	if i := sort.SearchStrings(c.nicks, n); i != len(c.nicks) && c.nicks[i] == n {
		c.nicks = append(c.nicks[:i], c.nicks[i+1:]...)
	}
}

func (c *channel) insertNick(n string) {
	i := sort.SearchStrings(c.nicks, n)
	if i < len(c.nicks) && c.nicks[i] == n {
		return
	}
	c.nicks = append(c.nicks, "")
	copy(c.nicks[i+1:], c.nicks[i:])
	c.nicks[i] = n
}
