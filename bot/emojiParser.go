package bot

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type Emoji struct {
	ID string
	Name string
}

func (e Emoji) asMessage() string {
	if len(e.ID) == 0 {
		// Simple emoji
		return e.Name
	}

	return fmt.Sprintf("<:%v:%v>", e.Name, e.ID)
}

func (e Emoji) asReaction() string {
	if len(e.ID) == 0 {
		// Simple emoji
		return e.Name
	}

	return fmt.Sprintf("%v:%v", e.Name, e.ID)
}

func findEmojiInMessage(s string) Emoji {
	if !strings.HasPrefix(s, "<") {
		// Simple emoji
		return Emoji{Name: s}
	}

	// Custom emoji, format: <:emoji_name:emoji_id>
	trimmed := strings.Trim(s,"<>")
	parts := strings.Split(trimmed,":")

	return Emoji{
		ID:   parts[2],
		Name: parts[1],
	}
}

func findEmojiIdentifierInReactionEmoji(emoji discordgo.Emoji) string {
	if len(emoji.ID) == 0 {
		// Simple emoji
		return emoji.Name
	}

	// Custom emoji
	return emoji.ID
}
