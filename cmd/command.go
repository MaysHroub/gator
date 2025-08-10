// Package cmd provides command-line argument parsing and executing utilities.
package cmd

type command struct {
	name string
	args []string 
}

func ParseCliArgs(args ...string) command {
	if len(args) == 0 {
		return command{}
	}
	return command{
		name: args[1],
		args: args[2:],
	}
}