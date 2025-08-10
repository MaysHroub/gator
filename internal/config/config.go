// Package config provides configuration structures and functions for the application.
package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const ConfigFileName = ".gatorconfig.json"

type ConfigService struct {
	configFilePath string
	cfg            config
}

func NewConfigService(filePath string) (*ConfigService, error) {
	cfg, err := readConfig(filePath)
	if err != nil {
		return nil, err
	}
	return &ConfigService{
		configFilePath: filePath,
		cfg: cfg,
	}, nil
} 

type config struct {
	DatabaseURL     string `json:"db_url"`
	CurrentUsername string `json:"current_user_name"`
}

func (cfgService *ConfigService) SetUser(username string) {
	cfgService.cfg.CurrentUsername = username
}

func (cfgService *ConfigService) GetConfig() config {
	return cfgService.cfg
}

func (cfgService *ConfigService) Save() error {
	return writeConfig(cfgService.cfg, cfgService.configFilePath)
}

func readConfig(path string) (config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return config{}, err
	}

	var cfg config
	if err = json.Unmarshal(data, &cfg); err != nil {
		return config{}, err
	}
	return cfg, nil
}

func writeConfig(cfg config, path string) error {
	data, err := json.Marshal(cfg)
	if err != nil {
		return err
	}
	if err = os.WriteFile(path, data, 0644); err != nil {
		return err
	}
	return nil
}

func GetConfigFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homeDir, ConfigFileName), nil
}
