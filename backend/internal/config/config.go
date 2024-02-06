package config

import (
	"chat/db"
	"chat/internal/cache"
	"chat/internal/router"
	"encoding/json"
	"os"
)

type Config struct {
	Server router.Config `json:"router"`
	Redis  cache.Config  `json:"cache"`
	DB     db.Config     `json:"db"`
}

func LoadConfiguration(configPath string) Config {
	configFile, err := os.Open(configPath)
	if err != nil {
		panic(err)
	}
	defer configFile.Close()

	var config Config
	decoder := json.NewDecoder(configFile)
	if err := decoder.Decode(&config); err != nil {
		panic(err)
	}

	return config
}
