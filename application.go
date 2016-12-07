package main

import (
	"fenrirc/config"
	"fenrirc/irc"
	"fenrirc/mondrian"
	"flag"
	"github.com/nsf/termbox-go"
	"io/ioutil"
	"time"
)

// Application pulls everything together.
type Application struct {
	*mondrian.Box

	current interface {
		mondrian.InteractiveWidget
		StatusProvider
		Handler
	}
	status *Status
	prompt *Prompt

	quit  bool
	cmd   chan func()
	event chan termbox.Event
	runUI func(func())

	frontends       []*Frontend
	currentFrontend int
	currentWidget   int
	clients         map[string]*irc.Client

	ticker *time.Ticker
}

// NewApplication is the constructor.
func NewApplication() *Application {
	ret := &Application{
		Box:           mondrian.NewBox(),
		current:       welcome,
		status:        NewStatus(),
		quit:          false,
		cmd:           make(chan func()),
		event:         make(chan termbox.Event),
		frontends:     []*Frontend{},
		clients:       make(map[string]*irc.Client),
		ticker:        time.NewTicker(1 * time.Second),
		currentWidget: -1,
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
	ret.status.Global = TimeStatusProvider{}

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
		a.Close()
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
	a.setCurrent(f.server)
	mondrian.Draw(a)

	go func() {
		if err := c.Connect(); err != nil {
			// TODO: log error
		} else {
			c.Run()
		}
	}()
}

func (a *Application) setCurrent(w interface {
	mondrian.InteractiveWidget
	Handler
	StatusProvider
}) {
	a.current.SetVisibility(false)
	a.current = w
	a.Box.Children[0] = w
	a.current.SetVisibility(true)
	a.Resize(a.Region)
	a.status.Current = a.current
}

func (a *Application) next() {
	if a.currentWidget+1 < len(a.frontends[a.currentFrontend].channels) {
		a.currentWidget++
		a.setCurrent(a.frontends[a.currentFrontend].channels[a.currentWidget])
	} else {
		a.currentFrontend = (a.currentFrontend + 1) % len(a.frontends)
		a.currentWidget = -1
		a.setCurrent(a.frontends[a.currentFrontend].server)
	}
}

func (a *Application) prev() {
	if a.currentWidget == 0 {
		a.currentWidget = -1
		a.setCurrent(a.frontends[a.currentFrontend].server)
	} else if a.currentWidget == -1 {
		a.currentFrontend = ((a.currentFrontend-1)%len(a.frontends) + len(a.frontends)) % len(a.frontends)
		a.currentWidget = len(a.frontends[a.currentFrontend].channels) - 1
		if a.currentWidget == -1 {
			a.setCurrent(a.frontends[a.currentFrontend].server)
		} else {
			a.setCurrent(a.frontends[a.currentFrontend].channels[a.currentWidget])
		}
	} else {
		a.currentWidget--
		a.setCurrent(a.frontends[a.currentFrontend].channels[a.currentWidget])
	}
}

// HandleKey handles user input.
func (a *Application) HandleKey(ev termbox.Event) {
	redraw := true
	switch ev.Key {
	case termbox.KeyCtrlQ:
		a.Close()
		return
	case termbox.KeyCtrlN:
		if len(a.frontends) == 0 {
			return
		}
		a.next()
	case termbox.KeyCtrlP:
		if len(a.frontends) == 0 {
			return
		}
		a.prev()
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

	old := time.Now()

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

		case t := <-a.ticker.C:
			if t.Minute() != old.Minute() {
				old = t
				mondrian.Draw(a.status)
			}
		}
	}
}

// Close closes all resources (IRC connections, ticker, ...) and quits the main loop.
func (a *Application) Close() {
	for _, client := range a.clients {
		if err := client.Close(); err != nil {
			// TODO: log error
		}
	}
	a.ticker.Stop()
	a.quit = true
}
