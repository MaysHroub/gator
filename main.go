package main

import (
	"fmt"
	"github/MaysHroub/gator/internal/config"
)

func main() {
	configFilePath, _ := config.GetConfigFilePath()
	cfg, err := config.ReadConfig(configFilePath)

	if err != nil {
		fmt.Printf("failed to get config instance: %v\n", err)
	}

	fmt.Println(cfg)
}

