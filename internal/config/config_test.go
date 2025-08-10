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
	// create a temporary file
	tempDir := t.TempDir()
	tempFilePath := filepath.Join(tempDir, "testread.json")

	// write test data to the file
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