package cmd

import (
	"github/MaysHroub/gator/internal/config"
	"github/MaysHroub/gator/internal/repository"
)

type state struct {
	cfg config.ConfigManager
	db  repository.Repository
}

func NewState(cfgMngr config.ConfigManager, db repository.Repository) *state {
	return &state{
		cfg:    cfgMngr,
		db: db,
	}
}
