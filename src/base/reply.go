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
