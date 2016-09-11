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

func newPrivate256(prefix, content string, toa time.Time) mondrian.Message {
	n := prefix
	if nickEnd := strings.Index(prefix, "!"); nickEnd != -1 {
		n = prefix[0:nickEnd]
	}
	return msg.Wrap(&private256{Nick: n, Content: content, ToA: toa})
}

func (p *private256) Draw(r *mondrian.Region) {
	r.Attr(termbox.ColorBlack|termbox.AttrBold, termbox.ColorDefault)
	r.LPrintf("[%02d:%02d] ", p.ToA.Hour(), p.ToA.Minute())
	r.AttrDefault()
	r.LPrintf("<")
	r.Attr(colorHash(p.Nick), termbox.ColorDefault)
	r.LPrintf("%s", p.Nick)
	r.AttrDefault()
	r.LPrintf("> ")
	r.Xbase = r.Cx
	r.Printf("%s", p.Content)
}
