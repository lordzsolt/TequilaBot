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
	config                      *Config
	currentState                state
	configurationStep           int
	roleReactionBeingConfigured *roleReaction
)

func Start(c *Config) error {
	config = c

	session, err := discordgo.New("Bot " + config.Token)
	if err != nil {
		return err
	}

	err = readReactions()
	if err != nil {
		return err
	}

	session.AddHandler(messageCreated)
	session.AddHandler(reactionAdd)
	session.AddHandler(reactionRemove)

	session.Identify.Intents =
		discordgo.IntentsGuildMessages |
		discordgo.IntentsGuildMessageReactions |
		discordgo.IntentsGuildMembers

	err = session.Open()
	if err != nil {
		return err
	}

	fmt.Println("Bot is running")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)

	<-sc

	return session.Close()
}

func isBotCommand(message *discordgo.MessageCreate, command string) bool {
	return strings.ToLower(message.Content) == command
}

func sendReply(session *discordgo.Session, message string) {
	fmt.Println("Sending reply: ", message)
	_, err := session.ChannelMessageSend(config.SetupChannelID, message)

	if err != nil {
		fmt.Println("Failed to send message: ", err.Error())
	}
}
