package cmd

import (
	"fmt"
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
	
	return nil
}
