// Package cmd provides command-line argument parsing and executing utilities.
package cmd

type Command struct {
	name string
	args []string
}

type commandInfo struct {
	name        string
	synopsis    string
	description string
	examples    []string
	author      string
	handler     func(st *State, cmd Command) error
}

func NewCommandInfo(name, synopsis, description, author string, examples []string, handler func(st *State, cmd Command) error) commandInfo {
	return commandInfo{
		name:        name,
		synopsis:    synopsis,
		description: description,
		examples:    examples,
		author:      author,
		handler:     handler,
	}
}

func ParseCliArgs(args ...string) Command {
	if len(args) == 0 {
		return Command{}
	}
	return Command{
		name: args[1],
		args: args[2:],
	}
}
