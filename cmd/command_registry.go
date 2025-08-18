package cmd

import "fmt"

type Commands struct {
	cmds map[string]commandInfo
}

func NewCommands() Commands {
	return Commands{
		cmds: make(map[string]commandInfo),
	}
}

func (c *Commands) Register(cmdName string, cmdInfo commandInfo) {
	c.cmds[cmdName] = cmdInfo
}

func (c *Commands) Run(st *State, cmd Command) error {
	cmdDetails, exists := c.cmds[cmd.name]
	if !exists {
		return fmt.Errorf("no such command exists: %s", cmd.name)
	}
	return cmdDetails.handler(st, cmd)
}
