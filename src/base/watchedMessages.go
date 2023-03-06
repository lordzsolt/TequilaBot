package base

import (
	"encoding/json"
	"fmt"
	"os"
)

const MessagesFileName = "./messages.json"

var (
	WatchedMessages = map[string]RoleReaction{}
)

func SaveMessages() error {
	bytes, err := json.MarshalIndent(WatchedMessages, "", "\t")

	if err != nil {
		fmt.Printf("Failed to marshal %v\n\n Error %v\n", WatchedMessages, err.Error())
		return err
	}

	err = os.WriteFile(MessagesFileName, bytes, 0644)

	if err != nil {
		fmt.Printf("Failed to write file: %v,\n", err)
		return err
	}

	return nil
}

func ReadMessages() error {
	fmt.Println("Reading reaction")

	_, err := os.Stat(MessagesFileName)
	if err != nil && os.IsNotExist(err) {
		fmt.Println("Reactions file doesn't exist yet.")
		return nil
	}

	file, err := os.ReadFile(MessagesFileName)
	if err != nil {
		fmt.Printf("Failed to read file: %v", err.Error())
		return err
	}

	err = json.Unmarshal(file, &WatchedMessages)
	if err != nil {
		fmt.Printf("Failed to unmarshal file: %v", err.Error())
		return err
	}

	return nil
}
