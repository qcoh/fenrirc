package event

var (
	// Events is a channel which receives functions supposed to run on a single goroutine.
	Events chan func()
)

func init() {
	Events = make(chan func())
}
