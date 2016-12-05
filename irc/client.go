package irc

import (
	"bufio"
	"crypto/tls"
	"fenrirc/config"
	"fenrirc/msg"
	"fmt"
	"net"
	"sort"
	"strings"
	"time"
)

// Client represents a connection to an IRC network.
type Client struct {
	conn     net.Conn
	conf     *config.Server
	Frontend Frontend

	channels map[string]Channel
	nicks    map[string][]string

	// run on ui goroutine
	runUI func(func())
}

// NewClient returns a client
func NewClient(conf *config.Server, runUI func(func())) *Client {
	return &Client{
		conf:     conf,
		runUI:    runUI,
		channels: make(map[string]Channel),
		nicks:    make(map[string][]string),
	}
}

// Connect connects client to IRC network.
func (c *Client) Connect() error {
	var err error
	if c.conf.SSL {
		c.conn, err = tls.Dial("tcp", c.conf.Host+":"+c.conf.Port, nil)
	} else {
		c.conn, err = net.Dial("tcp", c.conf.Host+":"+c.conf.Port)
	}
	if err != nil {
		return err
	}
	if c.conf.Pass != "" {
		c.Printf("PASS %s\r\n", c.conf.Pass)
	}
	c.Printf("NICK %s\r\n", c.conf.Nick)
	c.Printf("USER %s * * :%s\r\n", c.conf.User, c.conf.Real)
	return nil
}

// Write sends p to the server.
func (c *Client) Write(p []byte) (int, error) {
	// TODO: use a channel to make this somewhat async
	c.conn.SetWriteDeadline(time.Now().Add(50 * time.Millisecond))
	return c.conn.Write(p)
}

// Printf sends a formatted string to server.
func (c *Client) Printf(format string, a ...interface{}) {
	if _, err := fmt.Fprintf(c, format, a...); err != nil {
		c.logf("%s", err.Error())
	}
}

func (c *Client) logf(format string, a ...interface{}) {
	c.runUI(func() {
		c.Frontend.Server().Append(msg.NewDefault(&msg.Message{Raw: fmt.Sprintf(format, a...), ToA: time.Now()}))
	})
}

// Run spawns the read and write loops.
func (c *Client) Run() {
	go func() {
		scanner := bufio.NewScanner(c.conn)
		for scanner.Scan() {
			m, err := msg.Parse(scanner.Text())
			if err != nil {
				// handle error
				c.logf("Parsing error: %s", scanner.Text())
				continue
			}
			if m.Command == "PING" {
				c.Printf("PONG :%s\r\n", m.Trailing)
				continue
			}
			c.runUI(func() {
				c.handleMessage(m)
			})
		}
	}()
}

// returns the channel with name given by m.Params[n] if it exists, otherwise the server appender.
func (c *Client) appenderByParam(m *msg.Message, n int) Appender {
	if len(m.Params) <= n {
		return c.Frontend.Server()
	}
	if ch, ok := c.channels[m.Params[n]]; ok {
		return ch
	}
	return c.Frontend.Server()
}

func (c *Client) handleMessage(m *msg.Message) {
	switch m.Command {
	case "353", "RPL_NAMEREPLY":
		if len(m.Params) < 3 {
			c.Frontend.Server().Append(msg.NewDefault(m))
			return
		}
		c.nicks[m.Params[2]] = append(c.nicks[m.Params[2]], strings.Split(m.Trailing, " ")...)
	case "366", "RPL_ENDOFNAMES":
		if len(m.Params) < 2 {
			c.Frontend.Server().Append(msg.NewDefault(m))
			return
		}
		s := c.nicks[m.Params[1]]
		delete(c.nicks, m.Params[1])
		sort.Strings(s)
		c.appenderByParam(m, 1).Append(msg.NewNames(s, m.ToA))
	case "332", "RPL_TOPIC":
		a := c.appenderByParam(m, 1)
		if ch, ok := a.(Channel); ok {
			a.Append(msg.NewReplyTopic(m))
			ch.SetTopic(m.Trailing)
		} else {
			a.Append(msg.NewReplyTopic(m))
		}
	case "PRIVMSG":
		c.appenderByParam(m, 0).Append(msg.NewPrivate(m))
	case "JOIN":
		var name string
		if len(m.Params) > 0 {
			name = m.Params[0]
		} else {
			name = m.Trailing
		}
		var ch Channel
		var ok bool
		if ch, ok = c.channels[name]; !ok {
			ch = c.Frontend.NewChannel(name)
			c.channels[name] = ch
		}
		ch.Append(msg.NewJoin(m))
	default:
		c.Frontend.Server().Append(msg.NewDefault(m))
	}
}

// Close closes the client.
func (c *Client) Close() error {
	if c.conf.QuitMsg != "" {
		c.Printf("QUIT :%s\r\n", c.conf.QuitMsg)
	} else {
		c.Printf("QUIT\r\n")
	}
	err := c.conn.Close()

	return err
}
