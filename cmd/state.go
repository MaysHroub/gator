package cmd

import (
	"github/MaysHroub/gator/internal/config"
	"github/MaysHroub/gator/internal/repository"
)

type State struct {
	cfg config.ConfigManager
	db  repository.Repository
}

func NewState(cfgMngr config.ConfigManager, db repository.Repository) *State {
	return &State{
		cfg: cfgMngr,
		db:  db,
	}
}
