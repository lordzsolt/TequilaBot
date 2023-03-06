package bot

import (
	"fmt"
	"strconv"

	"github.com/bwmarrin/discordgo"

	"TequilaBot/src/base"
	"TequilaBot/src/config"
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

	rr, contains := base.WatchedMessages[reaction.MessageID]
	if !contains {
		return
	}

	emoji := reaction.Emoji.Name

	role, contains := rr.Reactions[emoji]
	if !contains {
		return
	}

	var err error
	if add {
		err = session.GuildMemberRoleAdd(config.Current.GuildID, reaction.UserID, role)
	} else {
		err = session.GuildMemberRoleRemove(config.Current.GuildID, reaction.UserID, role)
	}

	if err != nil {
		fmt.Printf("Failed to add: %v role %v to user %v because error %v\n", strconv.FormatBool(add), role, reaction.UserID, err.Error())
	} else {
		fmt.Printf("Successfully added: %v role %v to user %v\n", strconv.FormatBool(add), role, reaction.UserID)
	}
}
