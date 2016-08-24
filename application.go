package main

import (
	"fenrirc/config"
	"fenrirc/irc"
	"fenrirc/mondrian"
	"flag"
	"github.com/nsf/termbox-go"
)

// Application pulls everything together.
type Application struct {
	*mondrian.Box

	serverWindow *mondrian.MessageBuffer
	current      mondrian.InteractiveWidget
	status       *Status
	prompt       *Prompt

	quit  bool
	cmd   chan func()
	event chan termbox.Event
	runUI func(func())

	next Handler

	frontends map[string]*Frontend
	clients   map[string]*irc.Client
}

// NewApplication is the constructor.
func NewApplication() *Application {
	ret := &Application{
		Box:          mondrian.NewBox(),
		serverWindow: NewMessageBuffer(),
		status:       NewStatus(),
		quit:         false,
		cmd:          make(chan func()),
		event:        make(chan termbox.Event),
		frontends:    make(map[string]*Frontend),
		clients:      make(map[string]*irc.Client),
	}
	ret.current = ret.serverWindow
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
		if a.next != nil {
			a.next.Handle(cmd)
		}
	}
}

func (a *Application) connect(conf *config.Server) {
	if _, ok := a.clients[conf.Host]; ok {
		// reconnect?
		return
	}
	f := NewFrontend(a.serverWindow)
	a.frontends[conf.Host] = f
	c := irc.NewClient(f, conf, a.runUI)
	a.clients[conf.Host] = c

	go func() {
		if err := c.Connect(); err != nil {
			// TODO: log error
		} else {
			c.Run()
		}
	}()
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
