package flows

import (
	"github.com/bwmarrin/discordgo"
)

type Flow interface {
	Start(session *discordgo.Session)
	HandleMessage(session *discordgo.Session, message *discordgo.MessageCreate) (next Flow)
}
