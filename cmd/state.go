package cmd

import "github/MaysHroub/gator/internal/config"

type state struct {
	cfg config.ConfigManager
}

func NewState(cfgMngr config.ConfigManager) state {
	return state{
		cfg: cfgMngr,
	}
}