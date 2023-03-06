package flows

import (
	"github.com/bwmarrin/discordgo"

	"TequilaBot/src/base"
)

func NewConfigurationStartFlow() Flow {
	return &configurationStartFlow{}
}

const (
	prompt = "What would you like to do? Reply with 1-4\n" +
		"1. Add a new reaction message\n" +
		"2. Edit existing reaction message\n" +
		"3. Cleanup reaction messages\n" +
		"4. Nothing"
)

type configurationStartFlow struct {
}

func (f configurationStartFlow) Start(session *discordgo.Session) {
	base.SendReply(session, prompt)
}

func (f configurationStartFlow) HandleMessage(session *discordgo.Session, message *discordgo.MessageCreate) (next Flow) {
	switch message.Content[0:1] {
	case "1":
		return NewReactionsFlow()
	case "2":
		return NewEditingFlw()
	case "3":
		return NewCleanupFlow()
	case "4":
		return NewListeningFlow()
	default:
		base.SendUnknownMessage(session)
		return f
	}
}
