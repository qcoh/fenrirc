package irc

import (
	"fenrirc/cmd"
)

type channel struct {
	Channel
	*server
	name      string
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
