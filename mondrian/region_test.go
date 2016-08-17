package mondrian

import (
	"github.com/nsf/termbox-go"
	"testing"
)

func TestSetCell(t *testing.T) {
	SetMockUI(100, 100)
	r := &Region{Width: 100, Height: 100}

	r.SetCell(0, 1, '+', termbox.ColorDefault, termbox.ColorDefault)
	r.SetCell(10, 11, '-', termbox.ColorDefault, termbox.ColorDefault)

	if mockBuffer[0][1] != '+' {
		t.Errorf("SetCell(0,1,'+', ...) != '+'")
	}
	if mockBuffer[10][11] != '-' {
		t.Errorf("SetCell(10, 11, '-', ...) != '-'")
	}
}

func TestClear(t *testing.T) {
	SetMockUI(13, 27)
	r := &Region{Width: 13, Height: 27}

	for y := 0; y < r.Height; y++ {
		for x := 0; x < r.Width; x++ {
			r.SetCell(x, y, '+', r.Fg, r.Bg)
		}
	}
	r.Clear()

	for y := 0; y < r.Height; y++ {
		for x := 0; x < r.Width; x++ {
			if mockBuffer[x][y] != ' ' {
				t.Errorf("mockBuffer[%d][%d] != ' '", x, y)
			}
		}
	}
}

func TestComplicatedPrint(t *testing.T) {
	SetMockUI(10, 100)
	r := &Region{Width: 10, Height: 100}

	r.Printf("123456")
	if r.Cx != 6 {
		t.Errorf("r.Cx = %d, should be 6", r.Cx)
	}

	r.Printf("123456789ab")
	if r.Cy != 1 {
		t.Errorf("r.Cy = %d, should be 1 after newline", r.Cy)
	}
	if r.Cx != 1 {
		t.Errorf("r.Cx = %d, should be 1 after newline", r.Cx)
	}

	r.Move(0, 0)
	r.Printf("one two three")
	// should print like this:
	// >one two
	// >three
	if r.Cy != 1 || r.Cx != 5 {
		t.Errorf("r.Cx = %d, r.Cy = %d, should be (5,1)", r.Cx, r.Cy)
	}

	r.Move(0, 0)
	r.Xbase = 4
	r.Printf("one two three")
	// should print like this:
	// >....one
	// >....two
	// >....three
	if r.Cy != 2 || r.Cx != 9 {
		t.Errorf("r.Cx = %d, r.Cy = %d, should be (9,2)", r.Cx, r.Cy)
	}
}
