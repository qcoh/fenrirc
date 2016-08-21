package main

import (
	"fenrirc/mondrian"
	"github.com/nsf/termbox-go"
)

// Prompt wraps mondrian.Prompt.
type Prompt struct {
	*mondrian.Prompt
	Handler
}

// NewPrompt returns a prompt.
func NewPrompt(handler Handler) *Prompt {
	// TODO prompt bufsize
	return &Prompt{Prompt: mondrian.NewPrompt(512), Handler: handler}
}

// HandleKey handles user input.
func (p *Prompt) HandleKey(ev termbox.Event) {
	if ev.Key == termbox.KeyEnter {
		cmd, err := parse(p.Enter())
		if err != nil {
			// TODO: log error
		} else {
			p.Handle(cmd)
		}
	} else {
		p.Prompt.HandleKey(ev)
	}
}
