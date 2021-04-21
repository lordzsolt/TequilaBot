package bot

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type Config struct {
	Token     string `json:Token`
	BotPrefix string `json:BotPrefix`
	SetupChannelID string `json:SetupChannelID`
}

func ReadConfig() (*Config, error) {
	fmt.Println("Reading config")

	fmt.Println(os.Getwd())

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
