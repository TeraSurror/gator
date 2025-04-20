package main

import "fmt"

type command struct {
	Name string
	Args []string
}

type commands struct {
	commandFuncMap map[string]func(*state, command) error
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.commandFuncMap[name] = f
}

func (c *commands) run(s *state, cmd command) error {
	f, ok := c.commandFuncMap[cmd.Name]
	if !ok {
		return fmt.Errorf("command not supported: %s", cmd.Name)
	}

	return f(s, cmd)
}
