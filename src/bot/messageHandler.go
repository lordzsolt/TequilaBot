package bot

import (
	"github.com/bwmarrin/discordgo"

	"TequilaBot/src/base"
	"TequilaBot/src/config"
	"TequilaBot/src/flows"
)

var (
	currentFlow flows.Flow
)

func messageCreated(session *discordgo.Session, message *discordgo.MessageCreate) {
	if message.Author.ID == session.State.User.ID {
		return
	}
	if message.ChannelID != config.Current.SetupChannelID {
		return
	}

	var nextFlow flows.Flow
	if base.IsBotCommand(message, "abort") {
		nextFlow = flows.NewListeningFlow()
	} else {
		nextFlow = currentFlow.HandleMessage(session, message)
	}

	if nextFlow != currentFlow {
		nextFlow.Start(session)
		currentFlow = nextFlow
	}
}
