package main

import (
	"./mondrian"
)

// Prompt wraps mondrian.Prompt.
type Prompt struct {
	*mondrian.Prompt
}

// NewPrompt returns a prompt.
func NewPrompt() *Prompt {
	// TODO prompt bufsize
	return &Prompt{Prompt: mondrian.NewPrompt(512)}
}
