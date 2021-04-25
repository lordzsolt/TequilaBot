package main

import (
	"fmt"

	"TequilaBot/bot"
)
func main() {
	config, err := bot.ReadConfig()

	if err != nil {
		fmt.Printf("Failed to read config: %v\n", err.Error())
		return
	}

	err = bot.Start(config)

	if err != nil {
		fmt.Println(err.Error())
	}
}
