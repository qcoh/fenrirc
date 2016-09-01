package main

import (
	"fenrirc/mondrian"
	"fmt"
	"io"
)

type Server struct {
	*mondrian.MessageBuffer
	client io.Writer
}

func NewServer(client io.Writer) *Server {
	return &Server{MessageBuffer: NewMessageBuffer(), client: client}
}

func (s *Server) Handle(cmd *Command) {
	switch cmd.Command {
	case "WHOIS":
	case "JOIN":
		if len(cmd.Params) == 0 {
			return
		}
		fmt.Fprintf(s.client, "JOIN %s\r\n", cmd.Params[0])
		// TODO: write to messagebuffer
	}
}
