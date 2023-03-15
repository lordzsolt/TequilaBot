package base

import (
	"fmt"

	"github.com/bwmarrin/discordgo"

	"TequilaBot/src/config"
)

const (
	DiscordMessageLength = 2000
	UnknownMessage       = "Unknown message. Try answering the previous question again, or write **abort**"
)

func SendUnknownMessage(session *discordgo.Session) {
	SendReply(session, UnknownMessage)
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

func PostAllReactions(session *discordgo.Session) {
	var messageBody = "The following role messages are being tracked:\n"
	messageLength := len(messageBody)

	for _, message := range WatchedMessages {
		currentLine := fmt.Sprintf("**%v**: %v\n", message.Title, message.MessageID)
		currentLineLength := len(currentLine)
		if messageLength+currentLineLength > DiscordMessageLength {
			SendReply(session, messageBody)
			messageBody = ""
			messageLength = 0
		} else {
			messageBody += currentLine
			messageLength += currentLineLength
		}
	}

	if messageLength > 0 {
		SendReply(session, messageBody)
	}
}

func EditingExistingRoleReaction(session *discordgo.Session, reaction RoleReaction) {
	var edit = discordgo.MessageEdit{
		ID:      reaction.MessageID,
		Channel: reaction.ChannelID,
		Embed:   reaction.ToEmbed(),
	}

	_, err := session.ChannelMessageEditComplex(&edit)

	if err != nil {
		fmt.Println("Failed to edit message: ", err.Error())
	}
}
