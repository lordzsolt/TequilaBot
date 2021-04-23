package bot

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func messageCreated(session *discordgo.Session, message *discordgo.MessageCreate) {
	if message.Author.ID == session.State.User.ID {
		return
	}
	if message.ChannelID != config.SetupChannelID {
		return
	}

	fmt.Println(message.Content)

	switch currentState {
	case listening:
		if isBotCommand(message, config.BotPrefix) {
			startConfiguring(session)
		}
	case configuring:
		handleConfigurationMessage(session, message)
	}
}

func startConfiguring(session *discordgo.Session) {
	currentState = configuring
	configurationStep = 0
	messageBeingConfigured = &roleReactionMessage{
		reactions: map[string]string{},
	}
	sendReply(session, configurationMessages[configurationStep])
}

func handleConfigurationMessage(session *discordgo.Session, message *discordgo.MessageCreate) {
	if isBotCommand(message, config.SetupAbortPhrase) {
		currentState = listening
		return
	}

	switch configurationStep {
	case 0:
		channelID := strings.Trim(message.Content, "<#>")
		messageBeingConfigured.channelID = channelID
		reply := fmt.Sprintf("Alright. The channel is <#%s>", channelID)
		sendReply(session, reply)
	case 1:
		messageBeingConfigured.updateTitleAndDescription(message.Content)
	case 2:
		err := messageBeingConfigured.updateColor(message.Content)

		if err != nil {
			reply := fmt.Sprintf("Error converting color: %v\nTry again, or type `none` to leave message without color", err.Error())
			sendReply(session, reply)
			return
		}

		sendReply(session, "Perfect. This is how your message will look like:")
		messageBeingConfigured.messageID = sendCurrentMessage(session)
	case 3:
		if isBotCommand(message, "done") {
			// TODO: Add option to jump to steps and make corrections
			_= sendCurrentMessage(session)
			break
		}

		emoji, role := parseReaction(message.Content)

		if len(emoji) == 0 {
			fmt.Println("Could not parse emoji from: ", message)
			return
		}

		if len(role) == 0 {
			fmt.Println("Could not parse role from: ", message)
			return
		}

		// TODO: Make sure role exists

		messageBeingConfigured.reactions[emoji] = role
		messageBeingConfigured.messageID = updateCurrentMessage(session)
		err := session.MessageReactionAdd(config.SetupChannelID, messageBeingConfigured.messageID, emoji)
		if err != nil {
			fmt.Println("Error adding reaction: ", err.Error())
		}

		return
	case 4:
		if isBotCommand(message, "yes") {
			err := postFinalMessage(session)
			if err != nil {
				reply := fmt.Sprintf("Failed to post final message: %v\nType **yes** to try again.", err.Error())
				sendReply(session, reply)
			} else {
				currentState = listening
			}
		} else if isBotCommand(message, "no") {
			currentState = listening
		}

		return
	default:
		fmt.Println("Unknown configuration step: ", configurationStep)
		configurationStep = 0
		return
	}

	configurationStep += 1
	sendReply(session, configurationMessages[configurationStep])
}

func sendCurrentMessage(session *discordgo.Session) string {
	msg, err := session.ChannelMessageSendComplex(config.SetupChannelID, messageBeingConfigured.toDiscordMessage())

	if err != nil {
		fmt.Println("Failed to send message: ", err.Error())
	}

	if msg == nil {
		fmt.Println("For some reason, error is nil, but msg is also nil")
		return ""
	}

	return msg.ID
}

func updateCurrentMessage(session *discordgo.Session) string {
	msg, err := session.ChannelMessageEditEmbed(config.SetupChannelID, messageBeingConfigured.messageID, messageBeingConfigured.toEmbed())

	if err != nil {
		fmt.Println("Failed to edit message: ", err.Error())
	}

	if msg == nil {
		fmt.Println("For some reason, error is nil, but msg is also nil")
		return ""
	}

	return msg.ID
}

func parseReaction(message string) (emoji string, role string) {
	parts := strings.Fields(message)
	if len(parts) != 2 {
		fmt.Println("Could not create role from message: ", message)
		return "", ""
	}

	return parts[0], parts[1]
}

func postFinalMessage(session *discordgo.Session) error {
	msg, err := session.ChannelMessageSendComplex(messageBeingConfigured.channelID, messageBeingConfigured.toDiscordMessage())

	if err != nil {
		return err
	}

	for emoji, _ := range messageBeingConfigured.reactions {
		err := session.MessageReactionAdd(config.SetupChannelID, msg.ID, emoji)
		if err != nil {
			return err
		}
	}

	return nil
}
