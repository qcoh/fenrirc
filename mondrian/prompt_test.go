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
