package main

import (
	"./mondrian"
	"github.com/nsf/termbox-go"
)

// Application pulls everything together.
type Application struct {
	*mondrian.Box

	current mondrian.InteractiveWidget
	status  *Status
	prompt  *Prompt
}

// NewApplication is the constructor.
func NewApplication() *Application {
	ret := &Application{
		Box:     mondrian.NewBox(),
		current: NewMessageBuffer(),
		status:  NewStatus(),
		prompt:  NewPrompt(),
	}
	ret.Children = []mondrian.Widget{ret.current, ret.status, ret.prompt}
	ret.ResizeFunc = func(r *mondrian.Region) []*mondrian.Region {
		return []*mondrian.Region{
			{X: r.X, Y: r.Y, Width: r.Width, Height: r.Height - 2},
			{X: r.X, Y: r.Y + r.Height - 2, Width: r.Width, Height: 1},
			{X: r.X, Y: r.Y + r.Height - 1, Width: r.Width, Height: 1},
		}
	}
	return ret
}

// Run runs the application.
func (a *Application) Run() {
	event := make(chan termbox.Event)
	go func() {
		for {
			// TODO clean shutdown
			event <- mondrian.PollEvent()
		}
	}()

	w, h := mondrian.Size()
	a.SetVisibility(true)
	a.Resize(&mondrian.Region{Width: w, Height: h})
	mondrian.Draw(a)

mainloop:
	for {
		select {
		case ev := <-event:
			if ev.Type == termbox.EventResize {
				mondrian.Sync()
				a.Resize(&mondrian.Region{Width: ev.Width, Height: ev.Height})
				mondrian.Draw(a)
			}
			if ev.Ch != 0 {
				a.prompt.HandleKey(ev)
			}
			if ev.Key != 0 {
				a.prompt.HandleKey(ev)
				a.current.HandleKey(ev)
			}

			if ev.Key == termbox.KeyCtrlQ {
				break mainloop
			}
		}
	}
}
