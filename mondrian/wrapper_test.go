package mondrian

import (
	"github.com/nsf/termbox-go"
	"testing"
)

func TestWrappedFunctions(t *testing.T) {
	Init = func() error {
		return nil
	}
	Close = func() {
	}

	if err := Init(); err != nil {
		t.Error(err)
	}
	defer Close()
}

func TestSetMockUI(t *testing.T) {
	SetMockUI(60, 20)

	if err := Init(); err != nil {
		t.Error(err)
	}
	defer Close()

	Clear(termbox.ColorDefault, termbox.ColorDefault)
	Flush()

	SetCell(0, 10, 'Ö', termbox.ColorRed, termbox.ColorGreen)
	Flush()
}
