package flows

import (
	"github.com/bwmarrin/discordgo"
)

func NewEditingFlw() Flow {
	return editingFlow{}
}

type editingFlow struct {
}

func (e editingFlow) Start(session *discordgo.Session) {
	// TODO implement me
	panic("implement me")
}

func (e editingFlow) HandleMessage(session *discordgo.Session, message *discordgo.MessageCreate) (next Flow) {
	// TODO implement me
	panic("implement me")
}
