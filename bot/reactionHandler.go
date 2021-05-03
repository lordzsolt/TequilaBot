package bot

import (
	"fmt"
	"strconv"

	"github.com/bwmarrin/discordgo"
)

var (
	watchedMessages = map[string]roleReaction{}
)

func reactionAdd(session *discordgo.Session, reaction *discordgo.MessageReactionAdd) {
	handle(session, reaction.MessageReaction, true)
}

func reactionRemove(session *discordgo.Session, reaction *discordgo.MessageReactionRemove) {
	handle(session, reaction.MessageReaction, false)
}

func handle(session *discordgo.Session, reaction *discordgo.MessageReaction, add bool) {
	if reaction.UserID == session.State.User.ID {
		return
	}

	rr, contains := watchedMessages[reaction.MessageID]
	if !contains {
		return
	}

	emoji := findEmojiIdentifierInReactionEmoji(reaction.Emoji)

	role, contains := rr.Reactions[emoji]
	if !contains {
		return
	}

	var err error
	if add {
		err = session.GuildMemberRoleAdd(config.GuildID, reaction.UserID, role)
	} else {
		err = session.GuildMemberRoleRemove(config.GuildID, reaction.UserID, role)
	}

	if err != nil {
		fmt.Printf("Failed to add: %v role %v to user %v because error %v\n", strconv.FormatBool(add), role, reaction.UserID, err.Error())
	} else {
		fmt.Printf("Successfully added: %v role %v to user %v\n", strconv.FormatBool(add), role, reaction.UserID)
	}
}
