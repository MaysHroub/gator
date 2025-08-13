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

func TestStateCreationWithConfig_ValidCreation(t *testing.T) {
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
	name := "mays"
	mockConfig := config.MockConfigService{}
	mockConfig.On("SetUser", name).Return()
	mockConfig.On("Save").Return(nil)

	mockDB := repository.MockUserStore{}
	mockDB.On("GetUser", mock.Anything, name).Return(database.User{Name: name}, nil)

	st := NewState(&mockConfig, &mockDB)

	cmd := command{
		name: "login",
		args: []string{name},
	}

	err := HandleLogin(st, cmd)
	require.NoError(t, err)

	mockConfig.AssertCalled(t, "SetUser", name)
	mockConfig.AssertCalled(t, "Save")
	mockDB.AssertCalled(t, "GetUser", mock.Anything, name)
}

func TestLoginHandler_InvalidLogin_NoNameExists(t *testing.T) {
	name := "mays"
	
	mockDB := repository.MockUserStore{}
	mockDB.On("GetUser", mock.Anything, name).Return(database.User{}, nil)

	mockConfig := config.MockConfigService{}

	st := NewState(&mockConfig, &mockDB)
	cmd := command{
		name: "login",
		args: []string{name},
	}

	err := HandleLogin(st, cmd)
	require.Error(t, err)
	
	mockDB.AssertCalled(t, "GetUser", mock.Anything, name)
	mockConfig.AssertNotCalled(t, "SetUser", name)
	mockConfig.AssertNotCalled(t, "Save")
}

func TestRegisterHandler_ValidRegister(t *testing.T) {
	name := "mays"

	nameMatcher := mock.MatchedBy(func(p database.CreateUserParams) bool {
		return p.Name == name
	})

	mockDB := repository.MockUserStore{}
	mockDB.On(
		"GetUser",
		mock.Anything, // ctx
		name,
	).Return(database.User{}, nil)
	mockDB.
		On("CreateUser",
			mock.Anything, // ctx
			nameMatcher,
	).Return(database.User{}, nil)

	mockConfig := config.MockConfigService{}
	mockConfig.On("SetUser", name).Return()
	mockConfig.On("Save").Return(nil)

	st := NewState(&mockConfig, &mockDB)

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
	mockDB.AssertCalled(t,
		"GetUser",
		mock.Anything, // ctx
		name,
	)
	mockConfig.AssertCalled(t, "SetUser", name)
	mockConfig.AssertCalled(t, "Save")
}

func TestRegisterHandler_InvalidRegister_NameExists(t *testing.T) {
	name := "mays"

	mockDB := repository.MockUserStore{}
	mockDB.On(
		"GetUser",
		mock.Anything, // ctx
		name,
	).Return(database.User{Name: name}, nil)

	st := NewState(nil, &mockDB)

	cmd := command{
		name: "register",
		args: []string{name},
	}

	err := HandleRegister(st, cmd)
	require.Error(t, err)

	mockDB.AssertCalled(t,
		"GetUser",
		mock.Anything, // ctx
		name,
	)
	mockDB.AssertNotCalled(t, "CreateUser")
}
