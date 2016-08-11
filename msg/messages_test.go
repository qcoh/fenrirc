package msg

import (
	"testing"
)

func TestCachedMessage(t *testing.T) {
	m := &Cached{
		message:     &Simple{"just a test. what's up? nothing much... I guess I should finish my thesis instead of talking to myself"},
		heightCache: make(map[int]int),
	}

	for w := 20; w < 40; w++ {
		m.Height(w)
	}

	if len(m.heightCache) != 20 {
		t.Errorf("len(m.heightCache) = %d, should be 20", len(m.heightCache))
	}
}
