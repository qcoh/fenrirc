package irc

import (
	"fenrirc/cmd"
)

type server struct {
	Appender
	*Client
}

func (s *server) Handle(command *cmd.Command) {
	switch command.Method {
	case "JOIN":
		if len(command.Params) == 0 {
			return
		}
		s.server.Writef("JOIN %s\r\n", command.Params[0])
		// TODO: write to appender
	}
}
