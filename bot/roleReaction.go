package bot

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type roleReactionMessage struct {
	messageID   string
	channelID   string
	title       string
	description string
	color       int
	reactions   map[string]string
}

func (rrm roleReactionMessage) toEmbed() *discordgo.MessageEmbed {
	var embed discordgo.MessageEmbed
	embed.Title = rrm.title
	embed.Description = rrm.description

	if len(rrm.reactions) > 0 {
		embed.Description += "\n"
	}

	for emoji, role := range rrm.reactions {
		embed.Description += fmt.Sprintf("\n%v %v", emoji, role)
	}

	embed.Color = rrm.color
	return &embed
}

func (rrm roleReactionMessage) toDiscordMessage() *discordgo.MessageSend {
	return &discordgo.MessageSend{
		Embed: rrm.toEmbed(),
	}
}

func (rrm *roleReactionMessage) updateTitleAndDescription(message string) {
	parts := strings.Split(message, "|")
	messageBeingConfigured.title = parts[0]
	if len(parts) > 1 {
		messageBeingConfigured.description = parts[1]
	}
}

func (rrm *roleReactionMessage) updateColor(message string) error {
	if message == "none" {
		return nil
	}

	hex := strings.Trim(message, "#")
	value, err := strconv.ParseInt(hex, 16, 64)
	if err != nil {
		return err
	}

	messageBeingConfigured.color = int(value)
	return nil
}
