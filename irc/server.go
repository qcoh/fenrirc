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
		if _, ok := s.channels[command.Params[0]]; ok {
			// channel already exists, do nothing
			return
		}
		ch := &channel{server: s, name: command.Params[0], nicks: []string{}, nicksTemp: []string{}}
		ch.Channel = s.frontend.NewChannel(ch.name, ch)
		s.channels[ch.name] = ch
		s.server.Writef("JOIN %s\r\n", command.Params[0])
	}
}
