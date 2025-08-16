package cmd

import (
	"context"
	"errors"
	"fmt"
	"github/MaysHroub/gator/internal/database"
	"github/MaysHroub/gator/rss"
	"time"

	"github.com/google/uuid"
)

func HandleLogin(st *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("no enough args for %s; require username", cmd.name)
	}
	if !doesUserExist(st, cmd.args[0]) {
		return errors.New("user with given name doesn't exists")
	}
	st.cfg.SetCurrentUsername(cmd.args[0])
	st.cfg.Save()
	fmt.Println("current username got logged in")
	return nil
}

func HandleRegister(st *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("no enough args for %s; require username", cmd.name)
	}
	if doesUserExist(st, cmd.args[0]) {
		return errors.New("user with given name already exists")
	}
	ctx := context.Background()
	params := database.CreateUserParams{
		ID:        uuid.New(),
		Name:      cmd.args[0],
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	_, err := st.db.CreateUser(ctx, params)
	if err != nil {
		return err
	}
	st.cfg.SetCurrentUsername(cmd.args[0])
	st.cfg.Save()
	fmt.Println("current username got registered and logged in")
	return nil
}

func HandleResetUsers(st *state, cmd command) error {
	return st.db.DeleteAllUsers(context.Background())
}

func HandleListAllNames(st *state, cmd command) error {
	names, err := st.db.GetNamesOfAllUsers(context.Background())
	if err != nil {
		return err
	}
	for _, name := range names {
		if name == st.cfg.GetCurrentUsername() {
			fmt.Println(name + " (current)")
			continue
		}
		fmt.Println(name)
	}
	return nil
}

func HandleAgg(st *state, cmd command) error {
	url := "https://www.wagslane.dev/index.xml"
	rssFeed, err := rss.FetchFeed(url)
	if err != nil {
		return err
	}
	fmt.Println(rssFeed)
	return nil
}

func HandleAddFeed(st *state, cmd command, user database.User) error {
	if len(cmd.args) < 2 {
		return fmt.Errorf("no enough args for %s; require feed name and URL", cmd.name)
	}

	userID := user.ID
	feedID := uuid.New()

	createFeedParams := database.CreateFeedParams{
		ID:        feedID,
		Name:      cmd.args[0],
		Url:       cmd.args[1],
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    uuid.NullUUID{UUID: userID, Valid: true},
	}

	createFollowFeedParam := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    userID,
		FeedID:    feedID,
	}
	_, err := st.db.CreateFeed(context.Background(), createFeedParams)
	if err != nil {
		return err
	}

	_, err = st.db.CreateFeedFollow(context.Background(), createFollowFeedParam)

	return err
}

func HandleShowAllFeeds(st *state, cmd command) error {
	rows, err := st.db.GetAllFeeds(context.Background())
	if err != nil {
		return err
	}
	for _, row := range rows {
		fmt.Printf("Feed name: %s\nFeed URL: %s\nFeed creator name: %s\n------------------------>\n",
			row.Feedname, row.Url, row.Username)
	}
	return nil
}

func HandleFollowFeedByURL(st *state, cmd command, user database.User) error {
	if len(cmd.args) < 1 {
		return fmt.Errorf("no enough args for %s; require feed URL", cmd.name)
	}

	feedURL := cmd.args[0]

	feed, err := st.db.GetFeedByURL(context.Background(), feedURL)
	if err != nil {
		return err
	}

	_, err = st.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return err
	}

	fmt.Printf("User %s is following feed '%s' of ID %v\n", user.Name, feed.Name, feed.ID)

	return nil
}

func HandleShowAllFeedFollowsForUser(st *state, cmd command, user database.User) error {
	var username string
	if len(cmd.args) == 0 {
		username = user.Name
	} else {
		username = cmd.args[0]
	}

	res, err := st.db.GetFeedFollowsForUser(context.Background(), username)
	if err != nil {
		return err
	}

	fmt.Printf("Followed feeds of user %s:\n", username)

	for i, row := range res {
		fmt.Printf("%v. %s\n", i+1, row.FeedName)
	}

	return nil
}

func HandleUnfollowFeedByURL(st *state, cmd command, user database.User) error {
	if len(cmd.args) < 1 {
		return fmt.Errorf("no enough args for %s; require feed URL", cmd.name)
	}
	err := st.db.DeleteFeedFollowByUserAndURL(context.Background(), database.DeleteFeedFollowByUserAndURLParams{
		Name: user.Name,
		Url:  cmd.args[0],
	})
	return err
}

func doesUserExist(st *state, name string) bool {
	usr, err := st.db.GetUserByName(context.Background(), name)
	if err != nil {
		return false
	}
	if usr.Name == name {
		return true
	}
	return false
}

func MiddlewareLoggedIn(handler func(st *state, cmd command, user database.User) error) func(st *state, cmd command) error {
	return func(st *state, cmd command) error {
		user, err := st.db.GetUserByName(context.Background(), st.cfg.GetCurrentUsername())
		if err != nil {
			return err
		}
		return handler(st, cmd, user)
	}
}
