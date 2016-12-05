package main

import (
	"fenrirc/config"
	"fenrirc/mondrian"
	"fmt"
	"io"
)

// Server combines a messagebuffer with a handler.
type Server struct {
	*mondrian.MessageBuffer
	client io.Writer
	conf   *config.Server
}

// NewServer constructs a server.
func NewServer(conf *config.Server, client io.Writer) *Server {
	return &Server{MessageBuffer: NewMessageBuffer(), client: client, conf: conf}
}

// Handle handles user (prompt) input.
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

// Status provides server info.
func (s *Server) Status() string {
	// TODO: latency
	return fmt.Sprintf("%s [SSL: %t]", s.conf.Host, s.conf.SSL)
}
