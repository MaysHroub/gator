package main

import (
	"fmt"
	"os"
	"github/MaysHroub/gator/cmd"
	"github/MaysHroub/gator/internal/config"
)

func main() {
	filePath, _ := config.GetConfigFilePath()
	cfgService, err := config.NewConfigService(filePath)
	if err != nil {
		fmt.Printf("failed to instantiate config service: %v\n", err)
		return
	}
	st := cmd.NewState(cfgService)

	commands := cmd.NewCommands()
	commands.Register("login", cmd.HandleLogin)

	if len(os.Args) < 2 {
		fmt.Printf("no enough arguments; want at least 2, but got %v\n", len(os.Args))
		os.Exit(1)
	}

	cmnd := cmd.ParseCliArgs(os.Args...)

	err = commands.Run(&st, cmnd)
	if err != nil {
		fmt.Printf("command execution failed: %v\n", err)
		os.Exit(1)
	}
	os.Exit(0)
}
