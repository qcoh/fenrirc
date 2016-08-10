package mondrian

import (
	"github.com/mattn/go-runewidth"
	"github.com/nsf/termbox-go"
)

// Prompt is an input widget.
type Prompt struct {
	*Region
	Visible

	bufSize  int
	buffer   []rune
	startGap int
	endGap   int
}

// NewPrompt returns a prompt.
func NewPrompt(bufSize int) *Prompt {
	return &Prompt{
		Region:   defaultRegion,
		bufSize:  bufSize,
		buffer:   make([]rune, bufSize),
		startGap: 0,
		endGap:   bufSize - 1,
	}
}

// Resize resizes the prompt.
func (p *Prompt) Resize(r *Region) {
	p.Region = r
}

// Draw draws the prompt.
func (p *Prompt) Draw() {
	p.Clear()
	p.Move(0, 0)

	cursorPos := 0
	for _, ch := range p.buffer[0:p.startGap] {
		cursorPos += runewidth.RuneWidth(ch)
	}

	// Suppose for a moment that all runes in p.buffer have width = 1 and
	// split p.buffer in slices of length p.Width (padding the last one if
	// necessary). The following operation computes the starting index of
	// the slice containing p.buffer[cursorPos].
	// If p.buffer contains double-width characters, the operation produces
	// the correct index as well, provided we use 2 for the double-width
	// characters in the calculation of cursorPos (which we did).
	sliceStart := cursorPos / p.Width
	sliceStart *= p.Width

	// Draw the entire slice but shift out of the region for the correct start.
	p.Cx = -sliceStart

	for _, ch := range p.buffer[0:p.startGap] {
		p.SetCell(p.Cx, p.Cy, ch, p.Fg, p.Bg)
		p.Cx += runewidth.RuneWidth(ch)
	}
	SetCursor(p.X+p.Cx, p.Y+p.Cy)
	for _, ch := range p.buffer[p.endGap+1:] {
		// out of region
		if p.Cx >= p.Width {
			break
		}
		p.SetCell(p.Cx, p.Cy, ch, p.Fg, p.Bg)
		p.Cx += runewidth.RuneWidth(ch)
	}
}

// HandleKey handles user input and redraws the prompt if necessary.
func (p *Prompt) HandleKey(ev termbox.Event) {
	redraw := true

	if ev.Ch != 0 {
		p.insert(ev.Ch)
		goto draw
	}
	switch ev.Key {
	case termbox.KeySpace:
		p.insert(' ')
	case termbox.KeyArrowLeft:
		p.cursorLeft()
	case termbox.KeyArrowRight:
		p.cursorRight()
	case termbox.KeyBackspace, termbox.KeyBackspace2:
		if p.startGap > 0 {
			p.startGap--
		}
	case termbox.KeyDelete:
		if p.endGap < p.bufSize-1 {
			p.endGap++
		}
	case termbox.KeyCtrlA:
		for p.startGap > 0 {
			p.cursorLeft()
		}
	case termbox.KeyCtrlE:
		for p.endGap < p.bufSize-1 {
			p.cursorRight()
		}
	case termbox.KeyCtrlW:
		if p.startGap > 0 {
			p.startGap--
		}
		for p.startGap > 0 && p.buffer[p.startGap] != ' ' {
			p.startGap--
		}
	default:
		redraw = false
	}

draw:
	if redraw {
		Draw(p)
	}
}

// Enter returns the contents of the prompt, clears it and redraws it.
func (p *Prompt) Enter() string {
	ret := string(p.buffer[:p.startGap]) + string(p.buffer[p.endGap+1:])
	p.startGap = 0
	p.endGap = p.bufSize - 1
	Draw(p)
	return ret
}

func (p *Prompt) insert(ch rune) {
	if p.startGap < p.endGap {
		p.buffer[p.startGap] = ch
		p.startGap++
	}
}

func (p *Prompt) cursorLeft() {
	if p.startGap > 0 {
		p.buffer[p.endGap] = p.buffer[p.startGap-1]
		p.startGap--
		p.endGap--
	}
}

func (p *Prompt) cursorRight() {
	if p.endGap < p.bufSize-1 {
		p.buffer[p.startGap] = p.buffer[p.endGap+1]
		p.endGap++
		p.startGap++
	}
}
