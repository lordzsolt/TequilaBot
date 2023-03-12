package base

import (
	"fmt"

	"github.com/bwmarrin/discordgo"

	"TequilaBot/src/config"
)

const (
	unknownMessage = "Unknown message. Try answering the previous question again."
)

func SendUnknownMessage(session *discordgo.Session) {
	SendReply(session, unknownMessage)
}

func SendReply(session *discordgo.Session, message string) {
	fmt.Println("Sending reply: ", message)
	_, err := session.ChannelMessageSend(config.Current.SetupChannelID, message)

	if err != nil {
		fmt.Println("Failed to send message: ", err.Error())
	}
}

func SendReactionToChannel(session *discordgo.Session, reaction RoleReaction, channelID string, includeEmojis bool) (string, error) {
	var messageID string
	msg, err := session.ChannelMessageSendComplex(channelID, reaction.ToDiscordMessage())

	if err != nil {
		return messageID, err
	}
	messageID = msg.ID

	if !includeEmojis {
		return messageID, nil
	}

	for _, emoji := range reaction.Emojis {
		err := session.MessageReactionAdd(msg.ChannelID, msg.ID, emoji.AsReaction())
		if err != nil {
			return messageID, err
		}
	}

	return messageID, nil
}
