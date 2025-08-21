package cmd

import (
	"github/MaysHroub/gator/internal/config"
	"github/MaysHroub/gator/internal/database"
	"github/MaysHroub/gator/internal/repository"
	"path/filepath"
	"testing"

	"github.com/google/uuid"
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

	cfgService.SetCurrentUsername("mays-alreem")
	state := NewState(cfgService, nil)

	assert.Equal(t, cfgService, state.cfg)
}

func TestCommandsRegistryAndRun_ValidRegistryAndRun(t *testing.T) {
	cmdName := "login"
	called := false
	cmdHandler := func(st *State, cmd Command) error {
		called = true
		return nil
	}

	commands := NewCommands()
	commands.Register(cmdName, NewCommandInfo(cmdName, "", "", "", nil, cmdHandler))

	err := commands.Run(&State{}, Command{name: cmdName})
	require.NoError(t, err)

	assert.Equal(t, true, called)
}

func TestLoginHandler_ValidLogin(t *testing.T) {
	name := "mays"
	mockConfig := config.MockConfigService{}
	mockConfig.On("SetCurrentUsername", name).Return()
	mockConfig.On("GetCurrentUsername").Return("")
	mockConfig.On("Save").Return(nil)

	mockDB := repository.MockRepository{}
	mockDB.On("GetUserByName", mock.Anything, name).Return(database.User{Name: name}, nil)

	st := NewState(&mockConfig, &mockDB)

	cmd := Command{
		name: "login",
		args: []string{name},
	}

	err := HandleLogin(st, cmd)
	require.NoError(t, err)

	mockConfig.AssertCalled(t, "SetCurrentUsername", name)
	mockConfig.AssertCalled(t, "GetCurrentUsername")
	mockConfig.AssertCalled(t, "Save")
	mockDB.AssertCalled(t, "GetUserByName", mock.Anything, name)
}

func TestLoginHandler_InvalidLogin_NoNameExists(t *testing.T) {
	name := "mays"

	mockDB := repository.MockRepository{}
	mockDB.On("GetUserByName", mock.Anything, name).Return(database.User{}, nil)

	mockConfig := config.MockConfigService{}

	st := NewState(&mockConfig, &mockDB)
	cmd := Command{
		name: "login",
		args: []string{name},
	}

	err := HandleLogin(st, cmd)
	require.Error(t, err)

	mockDB.AssertCalled(t, "GetUserByName", mock.Anything, name)
	mockConfig.AssertNotCalled(t, "SetCurrentUsername", name)
	mockConfig.AssertNotCalled(t, "Save")
}

func TestRegisterHandler_ValidRegister(t *testing.T) {
	name := "mays"

	nameMatcher := mock.MatchedBy(func(p database.CreateUserParams) bool {
		return p.Name == name
	})

	mockDB := repository.MockRepository{}
	mockDB.On(
		"GetUserByName",
		mock.Anything, // ctx
		name,
	).Return(database.User{}, nil)
	mockDB.
		On("CreateUser",
			mock.Anything, // ctx
			nameMatcher,
		).Return(database.User{}, nil)

	mockConfig := config.MockConfigService{}
	mockConfig.On("SetCurrentUsername", name).Return()
	mockConfig.On("Save").Return(nil)

	st := NewState(&mockConfig, &mockDB)

	cmd := Command{
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
		"GetUserByName",
		mock.Anything, // ctx
		name,
	)
	mockConfig.AssertCalled(t, "SetCurrentUsername", name)
	mockConfig.AssertCalled(t, "Save")
}

func TestRegisterHandler_InvalidRegister_NameExists(t *testing.T) {
	name := "mays"

	mockDB := repository.MockRepository{}
	mockDB.On(
		"GetUserByName",
		mock.Anything, // ctx
		name,
	).Return(database.User{Name: name}, nil)

	st := NewState(nil, &mockDB)

	cmd := Command{
		name: "register",
		args: []string{name},
	}

	err := HandleRegister(st, cmd)
	require.Error(t, err)

	mockDB.AssertCalled(t,
		"GetUserByName",
		mock.Anything, // ctx
		name,
	)
	mockDB.AssertNotCalled(t, "CreateUser")
}

func TestResetCurrentUsernamesHandler(t *testing.T) {
	mockDB := repository.MockRepository{}
	mockDB.On("DeleteAllUsers", mock.Anything).Return(nil)

	st := NewState(nil, &mockDB)

	cmd := Command{
		name: "reset",
	}

	err := HandleResetUsers(st, cmd)
	require.NoError(t, err)

	mockDB.AssertCalled(t, "DeleteAllUsers", mock.Anything)
}

func TestListUsersNamesHandlers(t *testing.T) {
	mockDB := repository.MockRepository{}
	mockDB.On("GetNamesOfAllUsers", mock.Anything).Return([]string{"mays", "reem"}, nil)

	mockConfig := config.MockConfigService{}
	mockConfig.On("GetCurrentUsername").Return("mays")

	st := NewState(&mockConfig, &mockDB)
	cmd := Command{
		name: "users",
	}

	err := HandleListAllNames(st, cmd)
	require.NoError(t, err)

	mockDB.AssertCalled(t, "GetNamesOfAllUsers", mock.Anything)
}

func TestAddFeedHandler_ValidAddition(t *testing.T) {
	feedName := "feedname"
	feedURL := "https://example.com"
	user := database.User{
		ID:   uuid.New(),
		Name: "mays",
	}

	feedNameAndUserIDMatcher := mock.MatchedBy(func(p database.CreateFeedParams) bool {
		return p.Name == feedName && p.UserID.UUID == user.ID
	})
	userIDMatcher := mock.MatchedBy(func(p database.CreateFeedFollowParams) bool {
		return p.UserID == user.ID
	})

	mockDB := repository.MockRepository{}
	mockDB.On("CreateFeed", mock.Anything, feedNameAndUserIDMatcher).
		Return(database.Feed{Name: feedName, UserID: uuid.NullUUID{UUID: user.ID, Valid: true}}, nil)
	mockDB.On("CreateFeedFollow", mock.Anything, userIDMatcher).
		Return([]database.CreateFeedFollowRow{}, nil)

	st := NewState(nil, &mockDB)

	cmd := Command{
		name: "addfeed",
		args: []string{feedName, feedURL},
	}

	err := HandleAddFeed(st, cmd, user)
	require.NoError(t, err)

	mockDB.AssertCalled(t, "CreateFeed", mock.Anything, feedNameAndUserIDMatcher)
	mockDB.AssertCalled(t, "CreateFeedFollow", mock.Anything, userIDMatcher)
}

func TestShowAllFeedsHandler(t *testing.T) {
	rows := []database.GetAllFeedsRow{
		{
			Feedname: "feed1",
			Url:      "https://example1.com",
			Username: "user1",
		},
		{
			Feedname: "feed1",
			Url:      "https://example1.com",
			Username: "user2",
		},
		{
			Feedname: "feed2",
			Url:      "https://example1.com",
			Username: "user1",
		},
		{
			Feedname: "feed2",
			Url:      "https://example1.com",
			Username: "user3",
		},
	}

	mockDB := repository.MockRepository{}
	mockDB.On("GetAllFeeds", mock.Anything).Return(rows, nil)

	st := NewState(nil, &mockDB)
	cmd := Command{
		name: "feeds",
	}

	err := HandleShowAllFeeds(st, cmd)
	require.NoError(t, err)

	mockDB.AssertCalled(t, "GetAllFeeds", mock.Anything)
}

func TestFollowFeedHandler_ValidFollowing(t *testing.T) {
	feedID := uuid.New()
	feedname := "feed example"
	feedURL := "https://example.com"
	user := database.User{
		Name: "mays",
	}

	mockDB := repository.MockRepository{}
	mockDB.On("GetFeedByURL", mock.Anything, feedURL).Return(database.Feed{
		ID:   feedID,
		Name: feedname,
		Url:  feedURL,
	}, nil)

	createFeedFollowParamsMatcher := mock.MatchedBy(func(p database.CreateFeedFollowParams) bool {
		return p.FeedID == feedID && p.UserID == user.ID
	})
	mockDB.On("CreateFeedFollow", mock.Anything, createFeedFollowParamsMatcher).Return([]database.CreateFeedFollowRow{}, nil)

	st := NewState(nil, &mockDB)
	cmd := Command{name: "follow", args: []string{feedURL}}

	err := HandleFollowFeedByURL(st, cmd, user)
	require.NoError(t, err)

	mockDB.AssertCalled(t, "GetFeedByURL", mock.Anything, feedURL)
	mockDB.AssertCalled(t, "CreateFeedFollow", mock.Anything, createFeedFollowParamsMatcher)
}

func TestGetFeedFollowForUser_UsernameGivenInCmndArgs(t *testing.T) {
	user := database.User{
		Name: "mays",
	}
	feedID1, feedID2 := uuid.New(), uuid.New()
	feedname1, feedname2 := "feed example 1", "feed example 2"
	feedFollowRecords := []database.GetFeedFollowsForUserRow{
		{
			FeedID:   feedID1,
			FeedName: feedname1,
		},
		{
			FeedID:   feedID2,
			FeedName: feedname2,
		},
	}

	mockDB := repository.MockRepository{}
	mockDB.On("GetFeedFollowsForUser", mock.Anything, user.Name).Return(feedFollowRecords, nil)

	st := NewState(nil, &mockDB)
	cmd := Command{name: "following", args: []string{user.Name}}

	err := HandleShowAllFeedFollowsForUser(st, cmd, user)
	require.NoError(t, err)

	mockDB.AssertCalled(t, "GetFeedFollowsForUser", mock.Anything, user.Name)
}

func TestGetFeedFollowForUser_NoUsernameGivenInCmndArgs(t *testing.T) {
	user := database.User{
		Name: "mays",
	}
	feedID1, feedID2 := uuid.New(), uuid.New()
	feedname1, feedname2 := "feed example 1", "feed example 2"
	feedFollowRecords := []database.GetFeedFollowsForUserRow{
		{
			FeedID:   feedID1,
			FeedName: feedname1,
		},
		{
			FeedID:   feedID2,
			FeedName: feedname2,
		},
	}

	mockDB := repository.MockRepository{}
	mockDB.On("GetFeedFollowsForUser", mock.Anything, user.Name).Return(feedFollowRecords, nil)

	st := NewState(nil, &mockDB)
	cmd := Command{name: "following"}

	err := HandleShowAllFeedFollowsForUser(st, cmd, user)
	require.NoError(t, err)

	mockDB.AssertCalled(t, "GetFeedFollowsForUser", mock.Anything, user.Name)
}

func TestUnfollowFeedHandler(t *testing.T) {
	user := database.User{
		Name: "mays",
	}
	feedURL := "https://example.com"

	mockDB := repository.MockRepository{}
	params := database.DeleteFeedFollowByUserAndURLParams{
		Name: user.Name,
		Url:  feedURL,
	}
	mockDB.On("DeleteFeedFollowByUserAndURL", mock.Anything, params).Return(nil)

	st := NewState(nil, &mockDB)
	cmd := Command{name: "unfollow", args: []string{feedURL}}

	err := HandleUnfollowFeedByURL(st, cmd, user)
	require.NoError(t, err)

	mockDB.AssertCalled(t, "DeleteFeedFollowByUserAndURL", mock.Anything, params)
}

func TestBrowseHandler(t *testing.T) {
	user := database.User{
		Name: "mays",
	}
	cmd := Command{
		name: "browse",
	}
	defaultLimit := 2
	mockDB := repository.MockRepository{}
	mockDB.On(
		"GetPostsForUser",
		mock.Anything,
		database.GetPostsForUserParams{Name: user.Name, Limit: int32(defaultLimit)}).
		Return(
			[]database.Post{
				{
					Url: "https://example.com/post1",
				},
				{
					Url: "https://example.com/post2",
				},
			}, nil)

	st := NewState(nil, &mockDB)

	err := HandleBrowsePosts(st, cmd, user)
	require.NoError(t, err)

	mockDB.AssertCalled(t, "GetPostsForUser", mock.Anything, database.GetPostsForUserParams{Name: user.Name, Limit: int32(defaultLimit)})
}
