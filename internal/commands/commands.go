package commands

import (
	"fmt"

	"github.com/harshaditya-sharma/gator/internal/config"
)

type Command struct {
	Name string
	Args []string
}
type Commands struct {
	Handlers map[string]func(*config.State, Command) error
}

func (c *Commands) Register(name string, r func(*config.State, Command) error) {
	c.Handlers[name] = r
}

func (c *Commands) Run(s *config.State, cmd Command) error {
	handler, ok := c.Handlers[cmd.Name]
	if !ok {
		return fmt.Errorf("Unknown command: %s", cmd.Name)
	}
	return handler(s, cmd)
}
