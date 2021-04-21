package bot

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

var (
	config *Config
	currentState state
)

func Start(config *Config) error {
	session, err := discordgo.New("Bot " + config.Token)
	if err != nil {
		return err
	}

	session.AddHandler(messageCreated)

	session.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentsGuildMessageReactions

	err = session.Open()
	if err != nil {
		return err
	}

	fmt.Println("Bot is running")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)

	<- sc

	return session.Close()
}

func messageCreated(session *discordgo.Session, message *discordgo.MessageCreate) {
	if message.Author.ID == session.State.User.ID {
		return
	}

	if currentState == listening {
		if strings.HasPrefix(message.Content, config.BotPrefix) {
			currentState = 4
		}
	}

	fmt.Println("Received message: ", message.Content)
}
