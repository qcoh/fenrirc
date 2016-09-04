package mondrian

import (
	"testing"
)

type mockMsg struct {
	s string
}

func (mm *mockMsg) Draw(r *Region) {
	r.Printf("%s", mm.s)
}

func (mm *mockMsg) Height(width int) int {
	r := &Region{X: 1000, Y: 1000, Width: width, Height: 1000}
	mm.Draw(r)
	return r.Cy + 1
}

func TestMessageBufferDraw(t *testing.T) {
	SetMockUI(20, 100)
	m := NewMessageBuffer()
	m.SetVisibility(true)
	m.Resize(&Region{Width: 20, Height: 100})

	s0 := "this is a test"
	s1 := "123456789 123456789 next line"
	m.Append(&mockMsg{s0})
	m.Append(&mockMsg{s1})

	Draw(m)

	for k, ch := range s0 {
		if ch != mockBuffer[k][0] {
			t.Errorf("mockBuffer[%d][0] = %c != %c", k, mockBuffer[k][0], ch)
		}
	}

	for k, ch := range "next line" {
		if ch != mockBuffer[k][2] {
			t.Errorf("mockBuffer[%d][2] = %c != %c", k, mockBuffer[k][2], ch)
		}
	}
}

func TestMessageBufferScrolling(t *testing.T) {
	SetMockUI(10, 10)
	m := NewMessageBuffer()
	m.SetVisibility(true)
	m.Resize(&Region{Width: 10, Height: 10})

	for i := 0; i < 20; i++ {
		m.Append(&mockMsg{"foo"})
	}
	if m.totalHeight != 20 {
		t.Errorf("mockBuffer.totalHeight = %d != %d", m.totalHeight, 20)
	}
}
