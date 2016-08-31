package irc

import (
	"bufio"
	"crypto/tls"
	"fenrirc/config"
	"fenrirc/msg"
	"fmt"
	"net"
	"sync"
	"time"
)

// Client represents a connection to an IRC network.
type Client struct {
	conn     net.Conn
	conf     *config.Server
	frontend Frontend

	channels map[string]Channel

	// run on ui goroutine
	runUI func(func())

	// clean shutdown
	wg   sync.WaitGroup
	quit chan struct{}
}

// NewClient returns a client
func NewClient(frontend Frontend, conf *config.Server, runUI func(func())) *Client {
	return &Client{
		frontend: frontend,
		conf:     conf,
		runUI:    runUI,
		channels: make(map[string]Channel),
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
	n, err := c.conn.Write(p)
	if err != nil {
		c.logf("Timeout sending last message")
	}
	return n, err
}

// Printf sends a formatted string to server.
func (c *Client) Printf(format string, a ...interface{}) {
	fmt.Fprintf(c.conn, format, a...)
}

func (c *Client) logf(format string, a ...interface{}) {
	c.runUI(func() {
		c.frontend.Server().Append(msg.NewLog(fmt.Sprintf(format, a...), c.conf.Host, time.Now()))
	})
}

// Run spawns the read and write loops.
func (c *Client) Run() {
	go func() {
		c.wg.Add(1)
		defer c.wg.Done()

		scanner := bufio.NewScanner(c.conn)
		for scanner.Scan() {
			m, err := parse(scanner.Text())
			if err != nil {
				// handle error
				c.logf("Parsing error: %s", scanner.Text())
				continue
			}
			c.runUI(func() {
				c.handleMessage(m)
			})
		}
	}()
}

// returns the channel with name given by m.Params[n] if it exists, otherwise the server appender.
func (c *Client) appenderByParam(m *message, n int) Appender {
	if len(m.Params) <= n {
		return c.frontend.Server()
	}
	if ch, ok := c.channels[m.Params[n]]; ok {
		return ch
	}
	return c.frontend.Server()
}

func (c *Client) handleMessage(m *message) {
	switch m.Command {
	case "PRIVMSG":
		c.appenderByParam(m, 0).Append(msg.NewPrivate(m.Prefix, m.Trailing, m.ToA))
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
			ch = c.frontend.NewChannel(name)
			c.channels[name] = ch
		}
		ch.Append(msg.NewJoin(m.Prefix, name, m.ToA))
	case "PING":
		// writing to conn is thread safe. still might be better to do this in Run.
		c.Printf("PONG :%s\r\n", m.Trailing)
		fallthrough
	default:
		c.frontend.Server().Append(msg.NewDefault(c.conf.Host, m.Raw, m.ToA))
	}
}

// Close closes the client.
func (c *Client) Close() error {
	c.quit <- struct{}{}
	err := c.conn.Close()
	c.wg.Wait()

	return err
}
