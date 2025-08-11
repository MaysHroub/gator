package cmd

import (
	"context"
	"fmt"
	"github/MaysHroub/gator/internal/database"
	"time"

	"github.com/google/uuid"
)

func HandleLogin(st *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("no enough args for %s", cmd.name)
	}
	st.cfg.SetUser(cmd.args[0])
	st.cfg.Save()
	fmt.Println("current username got updated")
	return nil
}

func HandleRegister(st *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("no enough args for %s", cmd.name)
	}
	params := database.CreateUserParams{
        ID:        uuid.New(),
        Name:      cmd.args[0],
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    }
	_, err := st.db.CreateUser(context.Background(), params)
	return err
}
