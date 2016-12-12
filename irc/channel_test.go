package irc

import (
	"testing"
)

func TestHasNick(t *testing.T) {
	ch := &channel{nicks: []string{"aaa", "bbb", "ccc"}}
	if !ch.hasNick("bbb") {
		t.Errorf("hasNick(\"bbb\") == false")
	}

	if ch.hasNick("ddd") {
		t.Errorf("hasNick(\"ddd\") == true")
	}
}

func TestRemoveNick(t *testing.T) {
	ch := &channel{nicks: []string{}}
	ch.removeNick("aaa")

	ch = &channel{nicks: []string{"aaa"}}
	ch.removeNick("bbb")
	if ch.nicks[0] != "aaa" {
		t.Errorf("ch.nicks[0] != \"aaa\"")
	}

	ch.removeNick("aaa")
	if len(ch.nicks) != 0 {
		t.Errorf("ch.nicks != []string{}")
	}
}
