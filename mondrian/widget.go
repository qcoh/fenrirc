package mondrian

import (
	"github.com/nsf/termbox-go"
)

// Widget is the interface implemented by all ... widgets.
type Widget interface {
	Draw()
	Resize(*Region)
	IsVisible() bool
	SetVisibility(bool)
}

// InteractiveWidget is the interface implemented by all widgets which react to user input.
type InteractiveWidget interface {
	Widget
	HandleKey(termbox.Event)
}

// Draw draws arbitrarily many widgets. Flush is only called if at least one of them is visible.
func Draw(ws ...Widget) {
	flush := false
	for _, w := range ws {
		if w.IsVisible() {
			w.Draw()
			flush = true
		}
	}
	if flush {
		Flush()
	}
}
