// Package config provides configuration structures and functions for the application.
package config

type Config struct {
	DatabaseURL     string `json:"db_url"`
	CurrentUsername string `json:"current_user_name"`
}

func ReadConfig(path string) (Config, error) {
	return Config{}, nil
}