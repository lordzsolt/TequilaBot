package main

import (
	"fmt"

	"TequilaBot/src/bot"
	"TequilaBot/src/config"
)

func main() {
	config, err := config.ReadConfig()

	if err != nil {
		fmt.Printf("Failed to read config: %v\n", err.Error())
		return
	}

	err = bot.Start(*config)

	if err != nil {
		fmt.Println(err.Error())
	}
}
