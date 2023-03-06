package flows

import (
	"github.com/bwmarrin/discordgo"

	"TequilaBot/src/base"
	"TequilaBot/src/config"
)

func NewListeningFlow() Flow {
	return &listeningFlow{}
}

type listeningFlow struct {
}

func (f *listeningFlow) Start(session *discordgo.Session) {
	// Listening flow does not have any initial prompt, it only handles messages
}

func (f *listeningFlow) HandleMessage(session *discordgo.Session, message *discordgo.MessageCreate) (next Flow) {
	if base.IsBotCommand(message, config.Current.StartPhrase) {
		return NewConfigurationStartFlow()
	}
	return f
}
