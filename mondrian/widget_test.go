package mondrian

import (
	"testing"
)

type mockWidget struct {
	Widget
}

func (mockWidget) IsVisible() bool {
	return true
}

func (mockWidget) Draw() {
}

func TestDraw(t *testing.T) {
	SetMockUI(100, 100)

	a := mockWidget{}
	b := mockWidget{}

	Draw(a, b)
	Draw(a)
	Draw(b)
	Draw(a, b, b, a)
}
