package bot

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)


type Config struct {
	Token     string `json:Token`
	BotPrefix string `json:BotPrefix`
	SetupChannelID string `json:SetupChannelID`
	SetupAbortPhrase string `json:SetupAbortPhrase`
}

func ReadConfig() (*Config, error) {
	fmt.Println("Reading config")

	file, err := ioutil.ReadFile("./config.json")
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
