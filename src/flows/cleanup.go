package flows

import (
	"github.com/bwmarrin/discordgo"
)

func NewCleanupFlow() Flow {
	return cleanupFlow{}
}

type cleanupFlow struct {
}

func (c cleanupFlow) Start(session *discordgo.Session) {
	// TODO implement me
	panic("implement me")
}

func (c cleanupFlow) HandleMessage(session *discordgo.Session, message *discordgo.MessageCreate) (next Flow) {
	// TODO implement me
	panic("implement me")
}
