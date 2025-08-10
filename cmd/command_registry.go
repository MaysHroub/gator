package cmd

import "fmt"

type commands struct {
	cmds map[string]func(*state, command) error
}

func NewCommands() commands {
	mp := make(map[string]func(*state, command) error)
	return commands{
		cmds: mp,
	}
}

func (c *commands) Register(cmdName string, cmdHandler func(*state, command) error) {
	c.cmds[cmdName] = cmdHandler
}

func (c *commands) Run(st *state, cmd command) error {
	handler, exists := c.cmds[cmd.name]
	if !exists {
		return fmt.Errorf("no such command exists: %s", cmd.name)
	}
	return handler(st, cmd)
}
