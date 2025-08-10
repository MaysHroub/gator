package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreatingCommandWithParseCliArgs(t *testing.T) {
	args := []string{"login", "username"}
	cmd := ParseCLIArgs(args)

	assert.Equal(t, "login", cmd.name)
	assert.Equal(t, []string{"username"}, cmd.args)
}