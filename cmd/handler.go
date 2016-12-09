package cmd

// A Handler is the interface implemented by everything reacting to user (prompt) input.
type Handler interface {
	Handle(*Command)
}
