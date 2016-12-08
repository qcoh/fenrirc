package main

import (
	"fenrirc/config"
	"fenrirc/mondrian"
)

// Server combines a messagebuffer with a handler.
type Server struct {
	*mondrian.MessageBuffer
	conf *config.Server
}

// NewServer constructs a server.
func NewServer(conf *config.Server) *Server {
	return &Server{MessageBuffer: NewMessageBuffer(), conf: conf}
}
