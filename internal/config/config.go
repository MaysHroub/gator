// Package config provides configuration structures and functions for the application.
package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const ConfigFileName = ".gatorconfig.json"

type Config struct {
	DatabaseURL     string `json:"db_url"`
	CurrentUsername string `json:"current_user_name"`
}

func ReadConfig(path string) (Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}

	var config Config
	if err = json.Unmarshal(data, &config); err != nil {
		return Config{}, err
	}
	return config, nil
}

