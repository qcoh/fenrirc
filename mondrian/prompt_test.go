package mondrian

import (
	"github.com/nsf/termbox-go"
	"testing"
)

func TestPromptDraw(t *testing.T) {
	SetMockUI(100, 100)

	p := NewPrompt(512)
	p.SetVisibility(true)
	p.Resize(&Region{Width: 100, Height: 1})

	testString := "this is a test"
	for _, ch := range testString {
		p.HandleKey(termbox.Event{Ch: ch})
	}
	Draw(p)
	for k, ch := range testString {
		if mockBuffer[k][0] != ch {
			t.Errorf("mockBuffer[%d][0] = %c != %c", k, mockBuffer[k][0], ch)
		}
	}
}

func TestPromptCursor(t *testing.T) {
	SetMockUI(100, 100)
	SetCursor = func(x, y int) {
		mockBuffer[x][y] = '!'
	}
	p := NewPrompt(512)
	p.SetVisibility(true)
	p.Resize(&Region{Y: 50, Width: 100, Height: 1})

	p.insert('a')
	p.startGap--
	Draw(p)

	if mockBuffer[0][50] != '!' {
		t.Errorf("Cursor missing, mockBuffer[%d][%d] = %c != '!'", 0, 50, mockBuffer[0][50])
	}
}
