package cmd

import (
	"github/MaysHroub/gator/internal/config"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreatingCommandWithParseCliArgs_ValidParsing(t *testing.T) {
	args := []string{"login", "username"}
	cmd := ParseCliArgs(args...)

	assert.Equal(t, "login", cmd.name)
	assert.Equal(t, []string{"username"}, cmd.args)
}

func TestStateCreation_ValidCreation(t *testing.T) {
	path := filepath.Join(t.TempDir(), "statetest.json")
	cfgService, err := config.NewConfigService(path)
	require.NoError(t, err)
	
	cfgService.SetUser("mays-alreem")
	state := state {
		cfgService: &cfgService,
	}

	assert.Equal(t, cfgService, state.cfgService)
}