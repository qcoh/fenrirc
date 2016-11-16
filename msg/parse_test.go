package msg

import (
	"testing"
)

func compare(a, b *Message) bool {
	if len(a.Params) != len(b.Params) {
		return false
	}
	ret := a.Prefix == b.Prefix && a.Command == b.Command && a.Trailing == b.Trailing
	for k := range a.Params {
		ret = ret && a.Params[k] == b.Params[k]
	}
	return ret
}

func TestParse(t *testing.T) {
	in := []string{
		":Kevin!bncworld@I-Have.a.cool.vhost.com PRIVMSG #mIRC :I feel lucky today",
		"PING :something",
		":CalebDelnay!calebd@localhost MODE #mychannel -l",
		":CalebDelnay!calebd@localhost QUIT :Bye bye!",
		":Macha!~macha@unaffiliated/macha PRIVMSG #botwar :Test response",
	}
	out := []*Message{
		{Prefix: "Kevin!bncworld@I-Have.a.cool.vhost.com", Command: "PRIVMSG", Params: []string{"#mIRC"}, Trailing: "I feel lucky today"},
		{Prefix: "", Command: "PING", Params: []string{}, Trailing: "something"},
		{Prefix: "CalebDelnay!calebd@localhost", Command: "MODE", Params: []string{"#mychannel", "-l"}, Trailing: ""},
		{Prefix: "CalebDelnay!calebd@localhost", Command: "QUIT", Params: []string{}, Trailing: "Bye bye!"},
		{Prefix: "Macha!~macha@unaffiliated/macha", Command: "PRIVMSG", Params: []string{"#botwar"}, Trailing: "Test response"},
	}

	for k := range in {
		p, err := Parse(in[k])
		if err != nil {
			t.Error(err)
		}
		if !compare(p, out[k]) {
			t.Errorf("%+v != %+v\n", p, out[k])
		}
	}
}
