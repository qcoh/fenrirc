package main

import (
	"fenrirc/config"
	"fenrirc/irc"
	"fenrirc/mondrian"
	"github.com/nsf/termbox-go"
)

// Application pulls everything together.
type Application struct {
	*mondrian.Box

	current *mondrian.MessageBuffer
	status  *Status
	prompt  *Prompt

	quit  bool
	cmd   chan func()
	event chan termbox.Event
	runUI func(func())

	next Handler
}

// NewApplication is the constructor.
func NewApplication() *Application {
	ret := &Application{
		Box:     mondrian.NewBox(),
		current: NewMessageBuffer(),
		status:  NewStatus(),
		quit:    false,
		cmd:     make(chan func()),
		event:   make(chan termbox.Event),
	}
	ret.prompt = NewPrompt(ret)
	ret.Children = []mondrian.Widget{ret.current, ret.status, ret.prompt}
	ret.ResizeFunc = func(r *mondrian.Region) []*mondrian.Region {
		return []*mondrian.Region{
			{X: r.X, Y: r.Y, Width: r.Width, Height: r.Height - 2},
			{X: r.X, Y: r.Y + r.Height - 2, Width: r.Width, Height: 1},
			{X: r.X, Y: r.Y + r.Height - 1, Width: r.Width, Height: 1},
		}
	}

	// TODO: clean shutdown
	go func() {
		for {
			ret.event <- mondrian.PollEvent()
		}
	}()
	ret.runUI = func(f func()) {
		go func() {
			ret.cmd <- f
		}()
	}

	return ret
}

// Handle responds to "global" user commands.
func (a *Application) Handle(cmd *Command) {
	switch cmd.Command {
	case "QUIT":
		// TODO: clean shutdown
		a.quit = true
	default:
		if a.next != nil {
			a.next.Handle(cmd)
		}
	}
}

// Run runs the application.
func (a *Application) Run() {
	w, h := mondrian.Size()
	a.SetVisibility(true)
	a.Resize(&mondrian.Region{Width: w, Height: h})
	mondrian.Draw(a)

	conf := &config.Server{
		Host: "irc.freenode.net",
		Port: "6697",
		Nick: "qcoh_",
		User: "qcoh_",
		Real: "qcoh_",
		SSL:  true,
	}
	client := irc.NewClient(a, conf, a.runUI)
	// connect already uses cmd and blocks until cmd is emptied
	go func() {
		client.Connect()
		client.Run()
	}()

	for !a.quit {
		select {
		case ev := <-a.event:
			if ev.Type == termbox.EventResize {
				mondrian.Sync()
				a.Resize(&mondrian.Region{Width: ev.Width, Height: ev.Height})
				mondrian.Draw(a)
			} else if ev.Type == termbox.EventKey {
				if ev.Ch != 0 {
					a.prompt.HandleKey(ev)
				} else {
					a.prompt.HandleKey(ev)
					a.current.HandleKey(ev)
				}

				if ev.Key == termbox.KeyCtrlQ {
					a.quit = true
				}
			}
		case f := <-a.cmd:
			f()
		}
	}
}

// Server TODO: the following belong in their own frontend struct. I'm just testing if the previous work is correct.
func (a *Application) Server() irc.Appender {
	return a.current
}
