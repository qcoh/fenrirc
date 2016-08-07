package mondrian

import (
	"github.com/nsf/termbox-go"
)

var (
	// Init = termbox.Init
	Init = termbox.Init
	// Clear = termbox.Clear
	Clear = termbox.Clear
	// Close = termbox.Close
	Close = termbox.Close
	// Flush = termbox.Flush
	Flush = termbox.Flush
	// SetCell = termbox.SetCell
	SetCell = termbox.SetCell
	// SetCursor = termbox.SetCursor
	SetCursor = termbox.SetCursor
	// Size = termbox.Size
	Size = termbox.Size
	// Sync = termbox.Sync
	Sync = termbox.Sync
	// PollEvent = termbox.PollEvent
	PollEvent = termbox.PollEvent
	// SetOutputMode = termbox.SetOutputMode
	SetOutputMode = termbox.SetOutputMode
)

var (
	mockBuffer [][]rune
)

// SetMockUI overwrites all termbox functions with dummy functions.
func SetMockUI(w, h int) {
	mockBuffer = make([][]rune, w)
	for i := range mockBuffer {
		mockBuffer[i] = make([]rune, h)
	}

	Init = func() error {
		return nil
	}
	Clear = func(termbox.Attribute, termbox.Attribute) error {
		return nil
	}
	Close = func() {}
	Flush = func() error {
		return nil
	}
	SetCell = func(x int, y int, ch rune, _ termbox.Attribute, _ termbox.Attribute) {
		mockBuffer[x][y] = ch
	}
	SetCursor = func(x, y int) {}
	Size = func() (int, int) {
		return w, h
	}
	Sync = func() error {
		return nil
	}
	PollEvent = func() termbox.Event {
		return termbox.Event{}
	}
	SetOutputMode = func(termbox.OutputMode) termbox.OutputMode {
		return termbox.OutputMode(0)
	}
}
