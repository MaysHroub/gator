package cmd

import (
	"github/MaysHroub/gator/internal/config"
	"github/MaysHroub/gator/internal/database"
)

type state struct {
	cfg config.ConfigManager // an inteface is already a reference
	db *database.Queries
}

func NewState(cfgMngr config.ConfigManager, dbQrs *database.Queries) state {
	return state{
		cfg: cfgMngr,
		db: dbQrs,
	}
}