package cmd

import (
	"github/MaysHroub/gator/internal/config"
	"github/MaysHroub/gator/internal/database"
	"github/MaysHroub/gator/internal/repository"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCreatingCommandWithParseCliArgs_ValidParsing(t *testing.T) {
	cliArgs := []string{"rn", "login", "username"}
	cmd := ParseCliArgs(cliArgs...)

	assert.Equal(t, "login", cmd.name)
	assert.Equal(t, []string{"username"}, cmd.args)
}

func TestStateCreation_ValidCreation(t *testing.T) {
	path := filepath.Join(t.TempDir(), "test.json")
	config.WriteConfig(config.Config{}, path)

	cfgService, err := config.NewConfigService(path)
	require.NoError(t, err)

	cfgService.SetUser("mays-alreem")
	state := NewState(cfgService, nil)

	assert.Equal(t, cfgService, state.cfg)
}

func TestCommandsRegistryAndRun_ValidRegistryAndRun(t *testing.T) {
	cmdName := "login"
	called := false
	cmdHandler := func(st *state, cmd command) error {
		called = true
		return nil
	}

	commands := NewCommands()
	commands.Register(cmdName, cmdHandler)

	err := commands.Run(&state{}, command{name: cmdName})
	require.NoError(t, err)

	assert.Equal(t, true, called)
}

func TestLoginHandler_ValidLogin(t *testing.T) {
	mockConfig := config.MockConfigService{}
	mockConfig.On("SetUser", "mays-alreem").Return()
	mockConfig.On("Save").Return(nil)

	st := NewState(&mockConfig, nil)

	cmd := command{
		name: "login",
		args: []string{"mays-alreem"},
	}

	err := HandleLogin(st, cmd)
	require.NoError(t, err)

	mockConfig.AssertCalled(t, "SetUser", "mays-alreem")
	mockConfig.AssertCalled(t, "Save")
}

func TestRegisterHandler_ValidRegister(t *testing.T) {
	name := "mays"

	nameMatcher := mock.MatchedBy(func(p database.CreateUserParams) bool {
		return p.Name == name
	})

	mockDB := repository.MockUserStore{}
	mockDB.
		On("CreateUser",
        mock.Anything, // ctx
        nameMatcher,
    ).
    Return(database.User{}, nil)

	st := NewState(nil, &mockDB)

	cmd := command{
		name: "register",
		args: []string{name},
	}

	err := HandleRegister(st, cmd)
	require.NoError(t, err)

	mockDB.AssertCalled(t, 
		"CreateUser", 
		mock.Anything, // ctx
		nameMatcher,
	)
}