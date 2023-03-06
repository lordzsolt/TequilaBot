package config

import (
	"encoding/json"
	"fmt"
	"os"
)

var Current Config

type Config struct {
	Token          string `json:"token"`
	StartPhrase    string `json:"start_phrase"`
	GuildID        string `json:"guild_id"`
	SetupChannelID string `json:"setup_channel_id"`
}

func ReadConfig() (*Config, error) {
	fmt.Println("Reading config")

	file, err := os.ReadFile("./config.json")
	if err != nil {
		return nil, err
	}

	var config *Config
	err = json.Unmarshal(file, &config)

	if err != nil {
		return nil, err
	}

	return config, nil
}
