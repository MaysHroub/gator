package main

import (
	"database/sql"
	"fmt"
	"github/MaysHroub/gator/cmd"
	"github/MaysHroub/gator/internal/config"
	"github/MaysHroub/gator/internal/database"
	"os"

	_ "github.com/lib/pq"
) // importing the driver for a side effect; it won't be used

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
	cmnd := cmd.ParseCliArgs(os.Args...)

	RegisterCommands(commands, st, cmnd)

	err = commands.Run(st, cmnd)
	if err != nil {
		fmt.Printf("command execution failed: %v\n", err)
		os.Exit(1)
	}
}

func RegisterCommands(commands cmd.Commands, st *cmd.State, cmnd cmd.Command) {
	commands.Register("login", cmd.NewCommandInfo(
		"login",
		"login <username>",
		"Logs in a user to the CLI.",
		"mhroub",
		[]string{"login alice", "login bob"},
		cmd.HandleLogin,
	))

	commands.Register("register", cmd.NewCommandInfo(
		"register",
		"register <username>",
		"Registers a user and logs them in immediately.",
		"mhroub",
		[]string{"register alice"},
		cmd.HandleRegister,
	))

	commands.Register("users", cmd.NewCommandInfo(
		"users",
		"users",
		"Displays all registered users with '(current)' next to the logged-in user.",
		"mhroub",
		[]string{"users"},
		cmd.HandleListAllNames,
	))

	commands.Register("reset", cmd.NewCommandInfo(
		"reset",
		"reset",
		"Deletes all registered users.",
		"mhroub",
		[]string{"reset"},
		cmd.HandleResetUsers,
	))

	commands.Register("agg", cmd.NewCommandInfo(
		"agg",
		"agg",
		"Launches the feed aggregator in the background to fetch feeds and save posts in the database.",
		"mhroub",
		[]string{"agg"},
		cmd.HandleAgg,
	))

	commands.Register("addfeed", cmd.NewCommandInfo(
		"addfeed",
		"addfeed <feed-url>",
		"Adds a feed to the database. The user who adds it will be marked as the creator and automatically follow the feed.",
		"mhroub",
		[]string{"addfeed https://example.com/rss"},
		cmd.MiddlewareLoggedIn(cmd.HandleAddFeed),
	))

	commands.Register("feeds", cmd.NewCommandInfo(
		"feeds",
		"feeds",
		"Displays all feeds added to the database.",
		"mhroub",
		[]string{"feeds"},
		cmd.HandleShowAllFeeds,
	))

	commands.Register("follow", cmd.NewCommandInfo(
		"follow",
		"follow <feed-url>",
		"Follows a feed for the currently logged-in user.",
		"mhroub",
		[]string{"follow https://example.com/rss"},
		cmd.MiddlewareLoggedIn(cmd.HandleFollowFeedByURL),
	))

	commands.Register("unfollow", cmd.NewCommandInfo(
		"unfollow",
		"unfollow <feed-url>",
		"Unfollows a feed for the currently logged-in user.",
		"mhroub",
		[]string{"unfollow https://example.com/rss"},
		cmd.MiddlewareLoggedIn(cmd.HandleUnfollowFeedByURL),
	))

	commands.Register("following", cmd.NewCommandInfo(
		"following",
		"following [username]",
		"Displays all feeds followed by the given username. If omitted, shows feeds followed by the currently logged-in user.",
		"mhroub",
		[]string{"following", "following alice"},
		cmd.MiddlewareLoggedIn(cmd.HandleShowAllFeedFollowsForUser),
	))

	commands.Register("browse", cmd.NewCommandInfo(
		"browse",
		"browse [limit]",
		"Displays the latest posts. If limit is not provided, shows 2 posts. Each post shows title, description, and link.",
		"mhroub",
		[]string{"browse", "browse 5"},
		cmd.MiddlewareLoggedIn(cmd.HandleBrowsePosts),
	))

	commands.Register("man", cmd.NewCommandInfo(
		"man",
		"man [command]",
		"Displays the manual page of the specified command. If no command is provided, it explains how to use 'man'.",
		"mhroub",
		[]string{"man", "man users", "man feeds"},
		cmd.MiddlewareCommands(cmd.HandleShowManPage, commands),
	))

	commands.Register("cmnds", cmd.NewCommandInfo(
		"cmnds",
		"cmnds",
		"Displays all commands for this CLI tool.",
		"mhroub",
		[]string{"commands"},
		cmd.MiddlewareCommands(cmd.HandleShowAllCommands, commands),
	))
}

func getConfigService() (*config.ConfigService, error) {
	filePath, _ := config.GetConfigFilePath()
	cfgService, err := config.NewConfigService(filePath)
	if err != nil {
		return &config.ConfigService{}, err
	}
	return cfgService, nil
}
