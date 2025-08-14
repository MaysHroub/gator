package main

import (
	"database/sql"
	"fmt"
	"github/MaysHroub/gator/cmd"
	"github/MaysHroub/gator/internal/config"
	"github/MaysHroub/gator/internal/database"
	"os"

	_ "github.com/lib/pq"
) // importing the driver for a side effect; you won't use it

func main() {
	// check args length
	if len(os.Args) < 2 {
		fmt.Printf("no enough arguments; want at least 2, but got %v\n", len(os.Args))
		os.Exit(1)
	}

	// get configuration
	cfgService, err := getConfigService()
	if err != nil {
		fmt.Printf("failed to instantiate config service: %v\n", err)
		os.Exit(1)
	}

	// open db connection and create db queries
	db, err := sql.Open("postgres", cfgService.Cfg.DatabaseURL)
	if err != nil {
		fmt.Printf("failed to open db connection: %v\n", err)
		os.Exit(1)
	}
	dbQueries := database.New(db)

	st := cmd.NewState(cfgService, dbQueries)

	commands := cmd.NewCommands()
	commands.Register("login", cmd.HandleLogin)
	commands.Register("register", cmd.HandleRegister)
	commands.Register("reset", cmd.HandleResetUsers)
	commands.Register("users", cmd.HandleListAllNames)
	commands.Register("agg", cmd.HandleAgg)
	commands.Register("addfeed", cmd.HandleAddFeed)
	commands.Register("feeds", cmd.HandleShowAllFeeds)

	cmnd := cmd.ParseCliArgs(os.Args...)

	err = commands.Run(st, cmnd)
	if err != nil {
		fmt.Printf("command execution failed: %v\n", err)
		os.Exit(1)
	}
}

func getConfigService() (*config.ConfigService, error) {
	filePath, _ := config.GetConfigFilePath()
	cfgService, err := config.NewConfigService(filePath)
	if err != nil {
		return &config.ConfigService{}, err
	}
	return cfgService, nil
}

