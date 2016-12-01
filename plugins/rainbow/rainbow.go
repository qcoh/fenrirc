package rainbow

import (
	"fenrirc/mondrian"
	"fenrirc/msg"
	"github.com/nsf/termbox-go"
	"strings"
	"time"
)

func init() {
	mondrian.Init = mondrianInit256
	msg.NewPrivate = newPrivate256
}

func mondrianInit256() error {
	if err := termbox.Init(); err != nil {
		return err
	}
	termbox.SetOutputMode(termbox.Output256)
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	termbox.Flush()
	return nil
}

func colorHash(nick string) termbox.Attribute {
	const x = 127
	var ret uint = 0

	for _, v := range nick {
		ret = uint(v) + ret*x
		ret %= 216
	}
	return termbox.Attribute(ret + 0x10)
}

type private256 struct {
	Nick    string
	Content string
	ToA     time.Time
}

func newPrivate256(m *msg.Message) mondrian.Message {
	n := m.Prefix
	if nickEnd := strings.Index(m.Prefix, "!"); nickEnd != -1 {
		n = m.Prefix[0:nickEnd]
	}
	return msg.Wrap(&private256{Nick: n, Content: m.Trailing, ToA: m.ToA})
}

func (p *private256) Draw(r *mondrian.Region) {
	r.AttrDefault()
	r.LPrintf("[%02d:%02d] <", p.ToA.Hour(), p.ToA.Minute())
	r.Attr(colorHash(p.Nick), termbox.ColorDefault)
	r.LPrintf("%s", p.Nick)
	r.AttrDefault()
	r.LPrintf("> ")
	r.Xbase = r.Cx
	r.Printf("%s", p.Content)
}
