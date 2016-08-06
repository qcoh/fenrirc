package mondrian

import (
	"github.com/nsf/termbox-go"
	"testing"
)

func TestBoxDraw(t *testing.T) {
	SetMockUI(100, 100)
	var attrBuffer [100][100]termbox.Attribute
	SetCell = func(x, y int, ch rune, fg, bg termbox.Attribute) {
		attrBuffer[x][y] = bg
	}

	b := NewBox()
	b.ResizeFunc = func(r *Region) []*Region {
		return []*Region{
			{X: 0, Y: 0, Width: 100, Height: 50},
			{X: 0, Y: 50, Width: 100, Height: 50},
		}
	}

	d0 := &Dummy{}
	d1 := &Dummy{}

	b.Children = []Widget{d0, d1}

	b.Resize(&Region{Width: 100, Height: 100})
	b.SetVisibility(true)

	d0.Bg = termbox.ColorRed
	d1.Bg = termbox.ColorGreen

	Draw(b)

	for y := 0; y < 50; y++ {
		for x := 0; x < 100; x++ {
			if attrBuffer[x][y] != termbox.ColorRed {
				t.Errorf("attrBuffer[%d][%d] = %d != termbox.ColorRed", attrBuffer[x][y], x, y)
			}
		}
	}

	for y := 50; y < 100; y++ {
		for x := 0; x < 100; x++ {
			if attrBuffer[x][y] != termbox.ColorGreen {
				t.Errorf("attrBuffer[%d][%d] = %d != termbox.ColorGreen", attrBuffer[x][y], x, y)
			}
		}
	}
}
