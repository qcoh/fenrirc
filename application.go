package main

import (
	"fenrirc/config"
	"fenrirc/irc"
	"fenrirc/mondrian"
	"flag"
	"github.com/nsf/termbox-go"
	"io/ioutil"
)

// Application pulls everything together.
type Application struct {
	*mondrian.Box

	current mondrian.InteractiveWidget
	status  *Status
	prompt  *Prompt

	quit  bool
	cmd   chan func()
	event chan termbox.Event
	runUI func(func())

	frontends []*Frontend
	findex    int
	clients   map[string]*irc.Client
}

// NewApplication is the constructor.
func NewApplication() *Application {
	ret := &Application{
		Box:       mondrian.NewBox(),
		current:   firstMB,
		status:    NewStatus(),
		quit:      false,
		cmd:       make(chan func()),
		event:     make(chan termbox.Event),
		frontends: []*Frontend{},
		clients:   make(map[string]*irc.Client),
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
	case "CONNECT":
		// well, flag does the job for now
		fs := flag.NewFlagSet(cmd.Command, flag.ContinueOnError)
		fs.SetOutput(ioutil.Discard)
		conf := &config.Server{}
		fs.StringVar(&conf.Host, "Host", "", "")
		fs.StringVar(&conf.Port, "Port", "", "")
		fs.StringVar(&conf.Nick, "Nick", "", "")
		fs.StringVar(&conf.User, "User", "", "")
		fs.StringVar(&conf.Real, "Real", "", "")
		fs.StringVar(&conf.Pass, "Pass", "", "")
		fs.BoolVar(&conf.SSL, "SSL", false, "")
		if err := fs.Parse(cmd.Params); err != nil {
			// TODO: log error
		}
		a.connect(conf)

	default:
		if h, ok := a.current.(Handler); ok {
			h.Handle(cmd)
		}
	}
}

func (a *Application) connect(conf *config.Server) {
	if _, ok := a.clients[conf.Host]; ok {
		// TODO: reconnect?
		return
	}

	c := irc.NewClient(conf, a.runUI)
	a.clients[conf.Host] = c
	f := NewFrontend(conf, c)
	a.frontends = append(a.frontends, f)
	c.Frontend = f

	go func() {
		if err := c.Connect(); err != nil {
			// TODO: log error
		} else {
			c.Run()
		}
	}()
}

func (a *Application) setCurrent(w mondrian.InteractiveWidget) {
	a.current.SetVisibility(false)
	a.current = w
	a.Box.Children[0] = w
	a.current.SetVisibility(true)
	a.Resize(a.Region)
}

// HandleKey handles user input.
func (a *Application) HandleKey(ev termbox.Event) {
	var w mondrian.InteractiveWidget
	redraw := true
	switch ev.Key {
	case termbox.KeyCtrlQ:
		a.quit = true
		redraw = false
		return
	case termbox.KeyCtrlN:
		if len(a.frontends) == 0 {
			return
		}
		if w = a.frontends[a.findex].next(); w == nil {
			a.findex = (a.findex + 1) % len(a.frontends)
			w = a.frontends[a.findex].first()
		}
		a.setCurrent(w)
	case termbox.KeyCtrlP:
		if len(a.frontends) == 0 {
			return
		}
		if w = a.frontends[a.findex].prev(); w == nil {
			a.findex = ((a.findex-1)%len(a.frontends) + len(a.frontends)) % len(a.frontends)
			w = a.frontends[a.findex].last()
		}
		a.setCurrent(w)
	default:
		redraw = false
		a.prompt.HandleKey(ev)
		a.current.HandleKey(ev)
	}
	if redraw {
		a.Resize(a.Region)
		mondrian.Draw(a)
	}
}

// Run runs the application.
func (a *Application) Run() {
	w, h := mondrian.Size()
	a.SetVisibility(true)
	a.Resize(&mondrian.Region{Width: w, Height: h})
	mondrian.Draw(a)

	for !a.quit {
		select {
		case ev := <-a.event:
			if ev.Type == termbox.EventResize {
				mondrian.Sync()
				a.Resize(&mondrian.Region{Width: ev.Width, Height: ev.Height})
				mondrian.Draw(a)
			} else if ev.Type == termbox.EventKey {
				a.HandleKey(ev)
			}
		case f := <-a.cmd:
			f()
		}
	}
}
