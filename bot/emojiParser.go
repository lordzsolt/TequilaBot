package bot

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

func findEmojiIdentiferInMessage(s string) string {
	if !strings.HasPrefix(s, "<") {
		// Simple emoji
		return s
	}

	// Custom emoji, format: <:emoji_name:emoji_id>
	trimmer := strings.Trim(s,"<>")
	return strings.Split(trimmer, ":")[2]
}

func fincEmojiIdentifierInReactionEmoji(emoji discordgo.Emoji) string {
	if len(emoji.ID) == 0 {
		// Simple emoji
		return emoji.Name
	}

	// Custom emoji
	return emoji.ID
}
