package cmd

import (
	"context"
	"errors"
	"fmt"
	"github/MaysHroub/gator/internal/database"
	"github/MaysHroub/gator/rss"
	"strconv"
	"time"

	"github.com/google/uuid"
)

func HandleLogin(st *State, cmd Command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("no enough args for %s; require username", cmd.name)
	}
	if !doesUserExist(st, cmd.args[0]) {
		return errors.New("user with given name doesn't exists")
	}
	username := cmd.args[0]
	if st.cfg.GetCurrentUsername() == username {
		fmt.Printf("user %s is already logged in\n", username)
		return nil
	}
	st.cfg.SetCurrentUsername(username)
	st.cfg.Save()
	fmt.Printf("user %s got logged in\n", username)
	return nil
}

func HandleRegister(st *State, cmd Command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("no enough args for %s; require username", cmd.name)
	}
	if doesUserExist(st, cmd.args[0]) {
		return errors.New("user with given name already exists")
	}
	username := cmd.args[0]
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
	st.cfg.SetCurrentUsername(username)
	st.cfg.Save()
	fmt.Printf("user %s got registered and logged in\n", username)
	return nil
}

func HandleResetUsers(st *State, cmd Command) error {
	err := st.db.DeleteAllUsers(context.Background())
	if err != nil {
		return err
	}
	fmt.Println("all users were deleted")
	return nil
}

func HandleListAllNames(st *State, cmd Command) error {
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

func HandleAgg(st *State, cmd Command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("no enough args for %s; require time between requests", cmd.name)
	}
	timeBetweenReqs, err := time.ParseDuration(cmd.args[0])
	if err != nil {
		return err
	}
	fmt.Printf("Collecting feeds every %v\n", timeBetweenReqs)
	ticker := time.NewTicker(timeBetweenReqs)
	for ; ; <-ticker.C {
		err = rss.ScrapeFeeds(st.db)
		if err != nil {
			return err
		}
	}
}

func HandleAddFeed(st *State, cmd Command, user database.User) error {
	if len(cmd.args) < 2 {
		return fmt.Errorf("no enough args for %s; require feed name and URL", cmd.name)
	}

	userID := user.ID
	feedID := uuid.New()
	feedName := cmd.args[0]
	feedURL := cmd.args[1]

	createFeedParams := database.CreateFeedParams{
		ID:        feedID,
		Name:      feedName,
		Url:       feedURL,
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
	if err != nil {
		return err
	}
	fmt.Printf("feed %s is added\n", feedName)
	return nil
}

func HandleShowAllFeeds(st *State, cmd Command) error {
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

func HandleFollowFeedByURL(st *State, cmd Command, user database.User) error {
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

	fmt.Printf("you are now following feed '%s' of URL %v\n", feed.Name, feedURL)

	return nil
}

func HandleShowAllFeedFollowsForUser(st *State, cmd Command, user database.User) error {
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

	if len(res) == 0 {
		fmt.Printf("no follow feeds for user %s\n", username)
		return nil
	}

	fmt.Printf("Followed feeds of user %s:\n", username)

	for i, row := range res {
		fmt.Printf("%v. %s\n", i+1, row.FeedName)
	}

	return nil
}

func HandleUnfollowFeedByURL(st *State, cmd Command, user database.User) error {
	if len(cmd.args) < 1 {
		return fmt.Errorf("no enough args for %s; require feed URL", cmd.name)
	}
	feedURL := cmd.args[0]
	err := st.db.DeleteFeedFollowByUserAndURL(context.Background(), database.DeleteFeedFollowByUserAndURLParams{
		Name: user.Name,
		Url:  feedURL,
	})
	if err != nil {
		return err
	}
	fmt.Printf("you no longer follow feed of URL %s\n", feedURL)
	return nil
}

func HandleBrowsePosts(st *State, cmd Command, user database.User) error {
	const defaultLimit = 2
	limit := defaultLimit
	var err error
	if len(cmd.args) > 0 {
		limit, err = strconv.Atoi(cmd.args[0])
		if err != nil {
			return err
		}
	}
	posts, err := st.db.GetPostsForUser(context.Background(), database.GetPostsForUserParams{
		Name:  user.Name,
		Limit: int32(limit),
	})
	if err != nil {
		return err
	}
	for i, post := range posts {
		fmt.Printf("Post #%v\nTitle: %s\nDescription: %s\nLink: %s\n\n", i+1, post.Title, post.Description.String, post.Url)
	}
	return nil
}

func doesUserExist(st *State, name string) bool {
	usr, err := st.db.GetUserByName(context.Background(), name)
	if err != nil {
		return false
	}
	if usr.Name == name {
		return true
	}
	return false
}

func MiddlewareLoggedIn(handler func(st *State, cmd Command, user database.User) error) func(st *State, cmd Command) error {
	return func(st *State, cmd Command) error {
		user, err := st.db.GetUserByName(context.Background(), st.cfg.GetCurrentUsername())
		if err != nil {
			return err
		}
		return handler(st, cmd, user)
	}
}
