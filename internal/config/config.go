// Package config provides configuration structures and functions for the application.
package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const ConfigFileName = ".gatorconfig.json"

type ConfigManager interface {
	SetCurrentUsername(username string)
	GetCurrentUsername() string
	Save() error
}

type ConfigService struct {
	ConfigFilePath string
	Cfg            Config
}

func NewConfigService(filePath string) (*ConfigService, error) {
	cfg, err := ReadConfig(filePath)
	if err != nil {
		return nil, err
	}
	return &ConfigService{
		ConfigFilePath: filePath,
		Cfg:            cfg,
	}, nil
}

type Config struct {
	DatabaseURL     string `json:"db_url"`
	CurrentUsername string `json:"current_user_name"`
}

func (cfgService *ConfigService) SetCurrentUsername(username string) {
	cfgService.Cfg.CurrentUsername = username
}

func (cfgService *ConfigService) GetCurrentUsername() string {
	return cfgService.Cfg.CurrentUsername
}

func (cfgService *ConfigService) Save() error {
	return WriteConfig(cfgService.Cfg, cfgService.ConfigFilePath)
}

func ReadConfig(path string) (Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}

	var cfg Config
	if err = json.Unmarshal(data, &cfg); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

func WriteConfig(cfg Config, path string) error {
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
