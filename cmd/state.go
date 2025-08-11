package cmd

import (
	"github/MaysHroub/gator/internal/config"
	"github/MaysHroub/gator/internal/database"
	"github/MaysHroub/gator/internal/repository"
)

type state struct {
	cfg config.ConfigManager
	db repository.UserStore
}

func NewState(cfgMngr config.ConfigManager, dbQrs *database.Queries) state {
	return state{
		cfg: cfgMngr,
		db: dbQrs,
	}
}