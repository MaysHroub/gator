package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReadConfig_ReadsValidFile(t *testing.T) {
	tempDir := t.TempDir()
	tempFilePath := filepath.Join(tempDir, "testread.json")

	expected := Config{
		CurrentUsername: "mays",
		DatabaseURL:     "postgres://localhost:5432/mydb",
	}
	data, _ := json.Marshal(expected)
	err := os.WriteFile(tempFilePath, data, 0644)
	require.NoError(t, err)

	output, err := ReadConfig(tempFilePath)
	require.NoError(t, err)

	assert.Equal(t, expected, output)
}

func TestWriteConfig_WritesValidConfig(t *testing.T) {
	tempDir := t.TempDir()
	tempFilePath := filepath.Join(tempDir, "testread.json")

	cfg := Config{
		CurrentUsername: "mays",
		DatabaseURL:     "postgres://localhost:5432/mydb",
	}

	WriteConfig(cfg, tempFilePath)

	output, err := ReadConfig(tempFilePath)
	require.NoError(t, err)

	assert.Equal(t, cfg, output)
}

func TestConfigService_SetUserAndSave(t *testing.T) {
	// create a new service instance (empty)
	// sets the current username
	// save it
	// read the config and check
	tempDir := t.TempDir()
	tempFilePath := filepath.Join(tempDir, "testread.json")

	cfgService := NewConfigService(tempFilePath)

	expectedUsername := "mays-alreem"
	cfgService.SetUser(expectedUsername)

	err := cfgService.Save()
	require.NoError(t, err)

	cfg, err := ReadConfig(tempFilePath)
	require.NoError(t, err)

	assert.Equal(t, expectedUsername, cfg.CurrentUsername)
}