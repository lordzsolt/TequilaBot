package base

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

func IsBotCommand(message *discordgo.MessageCreate, s string) bool {
	return strings.ToLower(message.Content) == s
}
