package cmd

import "github/MaysHroub/gator/internal/config"

type state struct {
	cfgService *config.ConfigService
}

func NewState(cfgService *config.ConfigService) state {
	return state{
		cfgService: cfgService,
	}
}