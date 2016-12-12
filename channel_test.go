package main

import (
	"testing"
)

func TestHasNick(t *testing.T) {
	ch := &Channel{nicklist: &nicklist{nicks: []string{"aaa", "bbb", "ccc"}}}
	if !ch.HasNick("bbb") {
		t.Errorf("hasNick(\"bbb\") == false")
	}

	if ch.HasNick("ddd") {
		t.Errorf("hasNick(\"ddd\") == true")
	}
}

func TestRemoveNick(t *testing.T) {
	ch := &Channel{nicklist: &nicklist{nicks: []string{}}}
	ch.RemoveNick("aaa")

	ch = &Channel{nicklist: &nicklist{nicks: []string{"aaa"}}}
	ch.RemoveNick("bbb")
	if ch.nicks[0] != "aaa" {
		t.Errorf("ch.nicks[0] != \"aaa\"")
	}

	ch.RemoveNick("aaa")
	if len(ch.nicks) != 0 {
		t.Errorf("ch.nicks != []string{}")
	}
}

func TestInsertNick(t *testing.T) {
	ch := &Channel{nicklist: &nicklist{nicks: []string{}}}
	ch.InsertNick("aaa")
	if ch.nicks[0] != "aaa" {
		t.Errorf("ch.nicks[0] != \"aaa\"")
	}

	ch.InsertNick("bbb")
	if ch.nicks[1] != "bbb" {
		t.Errorf("ch.nicks[1] != \"bbb\"")
	}

	ch.InsertNick("aaa")
	if len(ch.nicks) != 2 {
		t.Errorf("len(ch.nicks) != 2")
	}
}

func TestSetNicks(t *testing.T) {
	ch := &Channel{nicklist: &nicklist{nicks: []string{}}}
	ch.SetNicks([]string{"bbb", "aaa"})

	if ch.nicks[0] != "aaa" || ch.nicks[1] != "bbb" {
		t.Errorf("ch.nicks[0] != \"aaa\" || ch.nicks[1] != \"bbb\"")
	}
}
