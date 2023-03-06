package bot

import (
	"encoding/json"
	"fmt"
	"os"
)

const ReactionFileName = "./messages.json"

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

func readReactions() error {
	fmt.Println("Reading reaction")

	_, err := os.Stat(ReactionFileName)
	if err != nil && os.IsNotExist(err) {
		fmt.Println("Reactions file doesn't exist yet.")
		return nil
	}

	file, err := os.ReadFile(ReactionFileName)
	if err != nil {
		fmt.Printf("Failed to read file: %v", err.Error())
		return err
	}

	err = json.Unmarshal(file, &watchedMessages)
	if err != nil {
		fmt.Printf("Failed to unmarshal file: %v", err.Error())
		return err
	}

	return nil
}

func saveReactions() error {
	bytes, err := json.MarshalIndent(watchedMessages, "", "\t")

	if err != nil {
		fmt.Printf("Failed to marshal %v\n\n Error %v\n", watchedMessages, err.Error())
		return err
	}

	err = os.WriteFile(ReactionFileName, bytes, 0644)

	if err != nil {
		fmt.Printf("Failed to write file: %v,\n", err)
		return err
	}

	return nil
}
