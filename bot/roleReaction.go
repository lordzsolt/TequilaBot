package bot

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type roleReaction struct {
	MessageID   string            `json:"message_id"`
	ChannelID   string            `json:"channel_id"`
	Title       string            `json:"title"`
	Description string            `json:"description"`
	Color       int               `json:"color"`
	Reactions   map[string]string `json:"reactions"`
	Emojis 		map[string]Emoji  `json:"emojis"`
}

func newRoleReactionMessage() *roleReaction {
	return &roleReaction{
		Reactions: map[string]string{},
		Emojis: map[string]Emoji{},
	}
}

func (rrm roleReaction) toEmbed() *discordgo.MessageEmbed {
	var embed discordgo.MessageEmbed
	embed.Title = rrm.Title
	embed.Description = rrm.Description

	if len(rrm.Reactions) > 0 {
		embed.Description += "\n"
	}

	for emojiName, role := range rrm.Reactions {
		emoji := rrm.Emojis[emojiName]
		embed.Description += fmt.Sprintf("\n%v <@&%v>", emoji.asMessage(), role)
	}

	embed.Color = rrm.Color
	return &embed
}

func (rrm roleReaction) toDiscordMessage() *discordgo.MessageSend {
	return &discordgo.MessageSend{
		Embed: rrm.toEmbed(),
	}
}

func (rrm *roleReaction) updateTitleAndDescription(message string) {
	parts := strings.Split(message, "|")
	roleReactionBeingConfigured.Title = parts[0]
	if len(parts) > 1 {
		roleReactionBeingConfigured.Description = parts[1]
	}
}

func (rrm *roleReaction) updateColor(message string) error {
	if message == "none" {
		return nil
	}

	hex := strings.Trim(message, "#")
	value, err := strconv.ParseInt(hex, 16, 64)
	if err != nil {
		return err
	}

	roleReactionBeingConfigured.Color = int(value)
	return nil
}
