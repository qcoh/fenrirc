package main

import (
	"fenrirc/cmd"
	"fenrirc/config"
	"fenrirc/mondrian"
)

// Server combines a messagebuffer with a handler.
type Server struct {
	*mondrian.MessageBuffer
	cmd.Handler
	conf *config.Server
}

// NewServer constructs a server.
func NewServer(conf *config.Server) *Server {
	return &Server{MessageBuffer: NewMessageBuffer(), conf: conf}
}
