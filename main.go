package main

import (
	"fmt"
	"github/MaysHroub/gator/internal/config"
)

func main() {
	filePath, _ := config.GetConfigFilePath()
	cfgService, err := config.NewConfigService(filePath)
	if err != nil {
		fmt.Printf("failed to instantiate config service: %v\n", err)
		return
	}

	cfgService.SetUser("mays-alreem")
	cfgService.Save()

	fmt.Println(cfgService.GetConfig())
}
