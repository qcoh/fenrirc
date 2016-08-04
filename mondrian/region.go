package mondrian

import (
	"bufio"
	"fmt"
	"github.com/mattn/go-runewidth"
	"github.com/nsf/termbox-go"
	"strings"
)

// Region contains the underlying position, size, cursor and color information for a widget.
type Region struct {
	// Position (top left corner)
	X, Y int
	// Size
	Width, Height int
	// Cursor position
	Cx, Cy int
	// Foreground and background attributes
	Fg, Bg termbox.Attribute
	// Baseline after a newline
	Xbase int
}

// SetCell prints a character to termbox' internal buffer.
func (r *Region) SetCell(x, y int, ch rune, fg, bg termbox.Attribute) {
	if 0 <= x && x < r.Width && 0 <= y && y < r.Height {
		SetCell(r.X+x, r.Y+y, ch, fg, bg)
	}
}

// Clear fills entire region with ' ' (space) with attributes r.Fg and r.Bg.
func (r *Region) Clear() {
	for y := 0; y < r.Height; y++ {
		for x := 0; x < r.Width; x++ {
			r.SetCell(x, y, ' ', r.Fg, r.Bg)
		}
	}
}

// Move moves the cursor.
func (r *Region) Move(x, y int) {
	r.Cx, r.Cy = x, y
}

// Attr sets the foreground and background attributes.
func (r *Region) Attr(fg, bg termbox.Attribute) {
	r.Fg, r.Bg = fg, bg
}

// AttrDefault sets the default attributes.
func (r *Region) AttrDefault() {
	r.Fg, r.Bg = termbox.ColorDefault, termbox.ColorDefault
}

func (r *Region) complicatedPrint(text string) {
	// TODO: benchmark this against simple strings.Split
	scanner := bufio.NewScanner(strings.NewReader(text))
	scanner.Split(bufio.ScanWords)

	for scanner.Scan() {
		word := scanner.Text()
		wordWidth := runewidth.StringWidth(word)

		// wordWidth exceeds maximum width, this print with linebreak.
		if wordWidth > r.Width-r.Xbase-1 {
			for _, ch := range word {
				// end of line reached: newline.
				if r.Cx%r.Width == 0 && r.Cx != 0 {
					r.Cx = r.Xbase
					r.Cy++
				}
				r.SetCell(r.Cx, r.Cy, ch, r.Fg, r.Bg)
				r.Cx += runewidth.RuneWidth(ch)
			}
			// whitespace at the end of the word (last whitespace is removed later!)
			r.SetCell(r.Cx, r.Cy, ' ', r.Fg, r.Bg)
			r.Cx++
			continue
		}
		// wordWidth exceeds remaining space in line but does not exceed maximum width, thus do a linebreak.
		if wordWidth > r.Width-r.Cx-1 {
			r.Cy++
			r.Cx = r.Xbase
		}

		// print entire word
		for _, ch := range word {
			r.SetCell(r.Cx, r.Cy, ch, r.Fg, r.Bg)
			r.Cx += runewidth.RuneWidth(ch)
		}

		// whitespace after the word
		r.SetCell(r.Cx, r.Cy, ' ', r.Fg, r.Bg)
		r.Cx++
	}
	// remove trailing whitespace
	r.Cx--
}

// Printf prints the string of a region. If the (visual) length exceeds r.Width-r.Xbase, the line is broken at a whitespace and printing continues on the next line. If the (visual) length of a word exceeds the maximum width, the linebreak is done in the middle of the word.
func (r *Region) Printf(format string, a ...interface{}) {
	r.complicatedPrint(fmt.Sprintf(format, a...))
}

func (r *Region) simplePrint(text string) {
	for _, ch := range text {
		r.SetCell(r.Cx, r.Cy, ch, r.Fg, r.Bg)
		r.Cx += runewidth.RuneWidth(ch)
		// not visible anymore
		if r.Cx > r.Width {
			break
		}
	}
}

// LPrintf prints the formatted string to a single line of a region. No linebreaks are performed.
func (r *Region) LPrintf(format string, a ...interface{}) {
	r.simplePrint(fmt.Sprintf(format, a...))
}
