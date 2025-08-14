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
		return fmt.Errorf("no enough args for %s", cmd.name)
	}
	if !doesUserExist(st, cmd.args[0]) {
		return errors.New("user with given name doesn't exists")
	}
	st.cfg.SetUser(cmd.args[0])
	st.cfg.Save()
	fmt.Println("current username got logged in")
	return nil
}

func HandleRegister(st *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("no enough args for %s", cmd.name)
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
	st.cfg.SetUser(cmd.args[0])
	st.cfg.Save()
	fmt.Println("current username got registered and logged in")
	return nil
}

func HandleResetUsers(st *state, cmd command) error {
	return st.db.DeleteAllUsers(context.Background())
}

func HandleListAllNames(st *state, cmd command) error {
	names, err := st.db.GetUsersNames(context.Background())
	if err != nil {
		return err
	}
	for _, name := range names {
		if name == st.cfg.GetUser() {
			fmt.Println(name + " (current)")
			continue
		}
		fmt.Println(name)
	}
	return nil
}

func HandleAgg(st *state, cmd command) error {
	url := "https://www.wagslane.dev/index.xml"
	rssFeed, err := rss.FetchFeed(context.Background(), url)
	if err != nil {
		return err
	}
	fmt.Println(rssFeed)
	return nil
}

func doesUserExist(st *state, name string) bool {
	usr, err := st.db.GetUser(context.Background(), name)
	if err != nil {
		return false
	}
	if usr.Name == name {
		return true
	}
	return false
}
