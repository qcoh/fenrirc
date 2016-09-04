package mondrian

import (
	"github.com/nsf/termbox-go"
)

// Message is an entry for a MessageBuffer.
type Message interface {
	Draw(*Region)
	Height(int) int
}

// MessageBuffer is an interactive widget showing a scrollable list of messages.
type MessageBuffer struct {
	*Region
	Visible

	messages    []Message
	totalHeight int
	scrollPos   int
	follow      bool
}

// NewMessageBuffer returns a MessageBuffer.
func NewMessageBuffer() *MessageBuffer {
	return &MessageBuffer{
		Region:   defaultRegion,
		messages: []Message{},
		follow:   true,
	}
}

// Resize resizes the MessageBuffer.
func (m *MessageBuffer) Resize(r *Region) {
	m.Region = r
	m.totalHeight = 0
	for _, msg := range m.messages {
		m.totalHeight += msg.Height(r.Width)
	}
}

// Draw draws the MessageBuffer.
func (m *MessageBuffer) Draw() {
	m.Move(0, 0)
	m.Clear()

	// scrolling beyond content/follow mode
	if m.scrollPos > m.totalHeight-m.Height || m.follow {
		m.scrollPos = m.totalHeight - m.Height
		m.follow = true
	}
	// if there are not enough lines to fill the screen, start at the top and
	// possibly leave some empty space.
	if m.scrollPos < 0 {
		m.scrollPos = 0
	}

	// find first message to draw
	m.Cy = -m.scrollPos
	startIndex := 0
	for m.Cy < 0 && len(m.messages[startIndex:]) > 0 {
		// found
		if m.Cy+m.messages[startIndex].Height(m.Width) > 0 {
			break
		}
		m.Cy += m.messages[startIndex].Height(m.Width)
		startIndex++
	}

	// draw
	for _, msg := range m.messages[startIndex:] {
		// outside of region, no need to draw anymore
		if m.Cy > m.Height {
			break
		}
		msg.Draw(m.Region)
		m.Cx = 0
		m.Cy++
	}
}

// HandleKey responds to user input.
func (m *MessageBuffer) HandleKey(ev termbox.Event) {
	redraw := true
	switch ev.Key {
	case termbox.KeyPgup:
		m.scrollPos -= m.Height / 2
		m.follow = false
	case termbox.KeyPgdn:
		m.scrollPos += m.Height / 2
		m.follow = false
	case termbox.KeyHome:
		m.scrollPos = 0
		m.follow = false
	case termbox.KeyEnd:
		m.follow = true
	default:
		redraw = false
	}
	if redraw {
		Draw(m)
	}
}

// Append appends a message.
func (m *MessageBuffer) Append(msg Message) {
	m.messages = append(m.messages, msg)
	m.totalHeight += msg.Height(m.Width)
	Draw(m)
}
