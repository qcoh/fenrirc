package irc

import (
	"../config"
	"bufio"
	"crypto/tls"
	"net"
	"sync"
	"time"
)

// Client represents a connection to an IRC network.
type Client struct {
	conn net.Conn
	conf *config.Server

	out chan []byte

	// run on ui goroutine
	cmd chan<- func()

	// clean shutdown
	wg   sync.WaitGroup
	quit chan struct{}
}

// NewClient returns a client
func NewClient(conf *config.Server, cmd chan<- func()) *Client {
	return &Client{
		conf: conf,
		cmd:  cmd,
	}
}

// Connect connects client to IRC network.
func (c *Client) Connect() error {
	var err error
	if c.conf.SLL {
		c.conn, err = tls.Dial("tcp", c.conf.Host+":"+c.conf.Port, nil)
	} else {
		c.conn, err = net.Dial("tcp", c.conf.Host+":"+c.conf.Port)
	}
	if err != nil {
		return err
	}
	if c.conf.Pass != "" {
		// TODO timeout
		fmt.Fprintf(c.conn, "PASS %s\r\n", c.conf.Pass)
	}
	// TODO timeout
	fmt.Fprintf(c.conn, "NICK %s\r\n", c.conf.Nick)
	fmt.Fprintf(c.conn, "USER %s * * :%s\r\n", c.conf.User, c.conf.Real)
	return nil
}

// Run spawns the read and write loops.
func (c *Client) Run() {
	go func() {
		c.wg.Add(1)
		defer wg.Done()

		scanner := bufio.NewScanner(c.conn)
		for scanner.Scan() {

		}
	}()

	go func() {
		c.wg.Add(1)
		defer c.wg.Done()

	loop:
		for {
			select {
			case <-quit:
				break loop
			case s := <-c.out:
				c.conn.SetWriteDeadline(time.No().Add(time.Second))
				if _, err := c.conn.Write(s); err != nil {
					// log error
				}
			}
		}
	}()
}

// Close closes the client.
func (c *Client) Close() error {
	c.quit <- struct{}{}
	err := c.conn.Close()
	c.wg.Wait()

	return err
}
