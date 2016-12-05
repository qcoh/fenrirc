package config

import (
	"github.com/BurntSushi/toml"
	"os"
)

// Global comprises the global configuration.
type Global struct {
	Servers map[string]Server
}

// Server comprises the server-specific configuration.
type Server struct {
	Host string
	Port string
	Nick string
	User string
	Real string
	Pass string
	SSL  bool

	QuitMsg string
}

// Get reads the configuration from the file `path` or returns an error.
func Get(path string) (*Global, error) {
	ret := &Global{}
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	if _, err := toml.DecodeReader(f, ret); err != nil {
		return nil, err
	}
	return ret, nil
}
