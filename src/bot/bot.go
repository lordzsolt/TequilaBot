package bot

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"

	"TequilaBot/src/base"
	"TequilaBot/src/config"
	"TequilaBot/src/flows"
)

func Start(c config.Config) error {
	config.Current = c

	session, err := discordgo.New("Bot " + c.Token)
	if err != nil {
		return err
	}

	err = base.ReadMessages()
	if err != nil {
		return err
	}

	currentFlow = flows.NewListeningFlow()

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
