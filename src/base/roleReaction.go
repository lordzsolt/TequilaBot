package base

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type RoleReaction struct {
	MessageID   string            `json:"message_id"`
	ChannelID   string            `json:"channel_id"`
	Title       string            `json:"title"`
	Description string            `json:"description"`
	Color       int               `json:"color"`
	Reactions   map[string]string `json:"reactions"`
	Emojis      map[string]Emoji  `json:"emojis"`
}

func NewRoleReaction() RoleReaction {
	return RoleReaction{
		Reactions: map[string]string{},
		Emojis:    map[string]Emoji{},
	}
}

func (rrm *RoleReaction) ToEmbed() *discordgo.MessageEmbed {
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

func (rrm *RoleReaction) ToDiscordMessage() *discordgo.MessageSend {
	return &discordgo.MessageSend{
		Embed: rrm.ToEmbed(),
	}
}

func (rrm *RoleReaction) UpdateTitleAndDescription(message string) {
	parts := strings.Split(message, "|")
	rrm.Title = parts[0]
	if len(parts) > 1 {
		rrm.Description = parts[1]
	}
}

func (rrm *RoleReaction) UpdateColor(message string) error {
	if message == "none" {
		return nil
	}

	hex := strings.Trim(message, "#")
	value, err := strconv.ParseInt(hex, 16, 64)
	if err != nil {
		return err
	}

	rrm.Color = int(value)
	return nil
}
