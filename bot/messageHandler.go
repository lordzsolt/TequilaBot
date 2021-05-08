package bot

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var configurationMessages = [...]string{
	"Step 1: Which channel would you like the message to be in?",
	"Step 2: How would you like your message to look like? Use a `|` to separate the title from the description, like so:\n" +
		"> `Title | Description`",
	"Step 3: Would you like your message to have a color? Respond with the hex color or **none**.",
	"Step 4: Finally, let's add roles. Add the emoji, then the name of the role.\n" +
		"Typing the same emoji a second time will override the previous role.\n" +
		// "To remove an emoji, type `remove :lemon:`\n" +
		"Type **done** when you are finished.\n" +
		"> Example: `:lemon: @Alcoholic`",
	"Good job! Above you can see how the final message looks like. Type **yes** to post it or **no** to discard it and start again",
}

func messageCreated(session *discordgo.Session, message *discordgo.MessageCreate) {
	if message.Author.ID == session.State.User.ID {
		return
	}
	if message.ChannelID != config.SetupChannelID {
		return
	}

	switch currentState {
	case listening:
		if isBotCommand(message, config.StartPhrase) {
			startConfiguring(session)
		}
	case configuring:
		handleConfigurationMessage(session, message)
	}
}

func startConfiguring(session *discordgo.Session) {
	currentState = configuring
	configurationStep = 0
	roleReactionBeingConfigured = newRoleReactionMessage()
	sendReply(session, configurationMessages[configurationStep])
}

func finishConfiguring() {
	watchedMessages[roleReactionBeingConfigured.MessageID] = *roleReactionBeingConfigured

	err := saveReactions()
	if err != nil {
		fmt.Printf("Failed to save reactions: %v", err.Error())
	}

	currentState = listening
	roleReactionBeingConfigured = nil
}

func handleConfigurationMessage(session *discordgo.Session, message *discordgo.MessageCreate) {
	if isBotCommand(message, "abort") {
		currentState = listening
		return
	}

	switch configurationStep {
	case 0:
		ChannelID := strings.Trim(message.Content, "<#>")
		roleReactionBeingConfigured.ChannelID = ChannelID
		reply := fmt.Sprintf("Alright. The channel is <#%s>", ChannelID)
		sendReply(session, reply)
	case 1:
		roleReactionBeingConfigured.updateTitleAndDescription(message.Content)
	case 2:
		err := roleReactionBeingConfigured.updateColor(message.Content)

		if err != nil {
			reply := fmt.Sprintf("Error converting color: %v\nTry again, or type `none` to leave message without color", err.Error())
			sendReply(session, reply)
			return
		}

		sendReply(session, "Perfect. This is how your message will look like:")
		postCurrentMessage(session)
	case 3:
		if isBotCommand(message, "done") {
			// TODO: Add option to jump to steps and make corrections
			postCurrentMessage(session)
			break
		}

		emoji, r := parseReaction(message.Content)

		// TODO: Make sure role exists

		roleReactionBeingConfigured.Reactions[emoji.Name] = r
		roleReactionBeingConfigured.Emojis[emoji.Name] = emoji
		roleReactionBeingConfigured.MessageID = updateCurrentMessage(session)

		err := session.MessageReactionAdd(config.SetupChannelID, roleReactionBeingConfigured.MessageID, emoji.asReaction())
		if err != nil {
			fmt.Println("Error adding reaction: ", err.Error())
		}

		return
	case 4:
		if isBotCommand(message, "yes") {
			err := postFinalMessage(session)
			if err != nil {
				reply := fmt.Sprintf("Failed to post final message: %v\nType **yes** to try again or **no** to abort and start from the beginning", err.Error())
				sendReply(session, reply)
			} else {
				finishConfiguring()
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

func updateCurrentMessage(session *discordgo.Session) string {
	var edit = discordgo.MessageEdit{
		ID: roleReactionBeingConfigured.MessageID,
		Channel: config.SetupChannelID,
		Embed: roleReactionBeingConfigured.toEmbed(),
	}

	msg, err := session.ChannelMessageEditComplex(&edit)

	if err != nil {
		fmt.Println("Failed to edit message: ", err.Error())
	}

	if msg == nil {
		fmt.Println("For some reason, error is nil, but msg is also nil")
		return ""
	}

	return msg.ID
}

func parseReaction(message string) (emoji Emoji, role string) {
	trimmed := strings.Trim(message, " ")
	parts := strings.Fields(trimmed)

	if len(parts) > 2 {
		fmt.Println("Warning: reaction / role message might have too many spaces")
	}

	emoji = findEmojiInMessage(parts[0])
	role = strings.Trim(parts[len(parts) - 1], "<@&>")
	return
}

func postCurrentMessage(session *discordgo.Session) {
	err := sendMessageToChannel(session, config.SetupChannelID)
	if err != nil {
		fmt.Printf("Failed to post current message: %v", err)
	}
}

func postFinalMessage(session *discordgo.Session) error {
	return sendMessageToChannel(session, roleReactionBeingConfigured.ChannelID)
}

func sendMessageToChannel(session *discordgo.Session, channelID string) error {
	msg, err := session.ChannelMessageSendComplex(channelID, roleReactionBeingConfigured.toDiscordMessage())

	if err != nil {
		return err
	}

	for _, emoji := range roleReactionBeingConfigured.Emojis {
		err := session.MessageReactionAdd(msg.ChannelID, msg.ID, emoji.asReaction())
		if err != nil {
			return err
		}
	}

	roleReactionBeingConfigured.MessageID = msg.ID

	return nil
}
