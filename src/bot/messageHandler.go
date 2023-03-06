package bot

import (
	"github.com/bwmarrin/discordgo"

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

	nextFlow := currentFlow.HandleMessage(session, message)

	if nextFlow != currentFlow {
		nextFlow.Start(session)
		currentFlow = nextFlow
	}
}
