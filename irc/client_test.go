package irc

import (
	"bytes"
	"net"
	"strings"
	"testing"
	"time"
)

type mockConn struct {
	net.Conn
	sr *strings.Reader
	bw bytes.Buffer
}

func (*mockConn) SetWriteDeadline(time.Time) error {
	return nil
}

func (m *mockConn) Read(b []byte) (int, error) {
	return m.sr.Read(b)
}

func (m *mockConn) Write(b []byte) (int, error) {
	return m.bw.Write(b)
}

func TestClientPingPong(t *testing.T) {
	input := "PING :1234\r\n"
	expected := "PONG :1234\r\n"

	mc := &mockConn{sr: strings.NewReader(input)}
	client := &Client{conn: mc}
	client.Run()

	resp := mc.bw.String()

	if resp != expected {
		t.Errorf("%s != %s\n", resp, expected)
	}
}
