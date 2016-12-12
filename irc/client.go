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

// A Client represents a connection to an IRC network.
type Client struct {
	conn     net.Conn
	conf     *config.Server
	frontend Frontend
	server   *server
	channels map[string]*channel
}

// NewClient constructs a client.
func NewClient(conf *config.Server, frontend Frontend) *Client {
	ret := &Client{
		conf:     conf,
		frontend: frontend,
		channels: make(map[string]*channel),
	}
	ret.server = &server{Client: ret}
	ret.server.Appender = ret.frontend.Server(ret.server)
	return ret
}

// Connect connects client to an IRC network.
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
		c.Writef("PASS %s\r\n", c.conf.Pass)
	}
	c.Writef("NICK %s\r\n", c.conf.Nick)
	c.Writef("USER %s * * :%s\r\n", c.conf.User, c.conf.Real)
	return nil
}

// Write sends p to the server.
func (c *Client) Write(p []byte) (int, error) {
	// TODO: make IO asynchronous using channels?
	c.conn.SetWriteDeadline(time.Now().Add(50 * time.Millisecond))
	return c.conn.Write(p)
}

// Writef sends a formatted string to the server.
func (c *Client) Writef(format string, a ...interface{}) {
	if _, err := fmt.Fprintf(c, format, a...); err != nil {
		c.logf("%s", err.Error())
	}
}

func (c *Client) logf(format string, a ...interface{}) {
	c.frontend.Sync(func() {
		c.server.Append(msg.NewLog(fmt.Sprintf(format, a...)))
	})
}

// Run reads from the server and dispatches its messages.
func (c *Client) Run() {
	scanner := bufio.NewScanner(c.conn)
	for scanner.Scan() {
		m, err := msg.Parse(scanner.Text())
		if err != nil {
			c.logf("Parsing error: %s", scanner.Text())
		} else if m.Command == "PING" {
			c.Writef("PONG :%s\r\n", m.Trailing)
		} else {
			c.frontend.Sync(func() {
				c.handleMessage(m)
			})
		}
	}
}

func (c *Client) channelByParam(m *msg.Message, n int) *channel {
	if len(m.Params) > n {
		if ch, ok := c.channels[m.Params[n]]; ok {
			return ch
		}
	}
	return nil
}

func (c *Client) handleMessage(m *msg.Message) {
	switch m.Command {
	case "353", "RPL_NAMEREPLY":
		if ch := c.channelByParam(m, 2); ch != nil {
			ch.nicksTemp = append(ch.nicksTemp, strings.Split(m.Trailing, " ")...)
		} else {
			c.logf("%s", m.Raw)
		}
	case "366", "RPL_ENDOFNAMES":
		if ch := c.channelByParam(m, 1); ch != nil {
			sort.Strings(ch.nicksTemp)
			ch.Append(msg.NewNames(ch.nicksTemp, m.ToA))
			ch.nicks = ch.nicksTemp
			ch.nicksTemp = []string{}
		} else {
			c.logf("%s", m.Raw)
		}
	case "332", "RPL_TOPIC":
		if ch := c.channelByParam(m, 1); ch != nil {
			ch.Append(msg.NewReplyTopic(m))
			ch.SetTopic(m.Trailing)
		} else {
			c.logf("%s", m.Raw)
		}
	case "333", "RPL_TOPICWHOTIME":
		if ch := c.channelByParam(m, 1); ch != nil {
			ch.Append(msg.NewReplyTopicWhoTime(m))
		} else {
			c.logf("%s", m.Raw)
		}
	case "NOTICE":
		c.server.Append(msg.NewNotice(m))
	//case "371", "RPL_INFO":
	case "372", "RPL_MOTD":
		c.server.Append(msg.NewMOTD(m))
	case "375", "RPL_MOTDSTART":
		c.server.Append(msg.NewMOTDStart(m))
	case "376", "RPL_ENDOFMOTD":
		c.server.Append(msg.NewEndOfMOTD(m))
	case "PRIVMSG":
		if ch := c.channelByParam(m, 0); ch != nil {
			ch.Append(msg.NewPrivate(m))
		} else {
			c.logf("%s", m.Raw)
		}
	case "JOIN":
		// Different servers have the channel name in different places.
		name := m.Trailing
		if len(m.Params) > 0 {
			name = m.Params[0]
		}
		if ch, ok := c.channels[name]; ok {
			ch.Append(msg.NewJoin(m))
			if n, _, ok := nickHost(m.Prefix); ok {
				ch.insertNick(n)
			}
		} else {
			c.logf("%s", m.Raw)
		}
	case "QUIT":
		if n, _, ok := nickHost(m.Prefix); ok {
			for _, ch := range c.channels {
				if ch.hasNick(n) {
					ch.Append(msg.NewQuit(m))
					ch.removeNick(n)
				}
			}
		} else {
			c.logf("%s", m.Raw)
		}
	case "NICK":
		if n, _, ok := nickHost(m.Prefix); ok && m.Trailing != "" {
			for _, ch := range c.channels {
				if ch.hasNick(n) {
					ch.removeNick(n)
					ch.insertNick(m.Trailing)
					ch.Append(msg.NewNick(m))
				}
			}
		} else {
			c.logf("%s", m.Raw)
		}

	default:
		c.server.Append(msg.NewDefault(m))
	}
}

// Close closes the IRC connection.
func (c *Client) Close() error {
	if c.conf.QuitMsg != "" {
		c.Writef("QUIT :%s\r\n", c.conf.QuitMsg)
	} else {
		c.Writef("QUIT\r\n", c.conf.QuitMsg)
	}
	return c.conn.Close()
}

func nickHost(s string) (string, string, bool) {
	if nickEnd := strings.Index(s, "!"); nickEnd != -1 {
		// TODO: verify nickEnd+2 < len
		return s[0:nickEnd], s[nickEnd+2:], true
	}
	return s, "", false
}
