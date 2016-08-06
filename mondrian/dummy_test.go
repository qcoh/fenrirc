package mondrian

import (
	"github.com/nsf/termbox-go"
	"testing"
)

func TestDummyDraw(t *testing.T) {
	SetMockUI(100, 100)
	var attrBuffer [100][100]termbox.Attribute
	SetCell = func(x, y int, ch rune, fg, bg termbox.Attribute) {
		attrBuffer[x][y] = bg
	}

	d := &Dummy{
		Region{Width: 100, Height: 100, Bg: termbox.ColorGreen},
		Visible{true},
	}
	Draw(d)

	for y := 0; y < 100; y++ {
		for x := 0; x < 100; x++ {
			if attrBuffer[x][y] != termbox.ColorGreen {
				t.Errorf("attrBuffer[%d][%d] != termbox.ColorGreen", x, y)
			}
		}
	}
}
