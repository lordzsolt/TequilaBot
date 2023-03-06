package flows

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"

	"TequilaBot/src/base"
	"TequilaBot/src/config"
)

func NewReactionsFlow() Flow {
	return &reactionsFlow{}
}

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

type reactionsFlow struct {
	step    int
	message base.RoleReaction
}

func (f *reactionsFlow) Start(session *discordgo.Session) {
	f.step = 0
	f.message = base.NewRoleReaction()
	base.SendReply(session, configurationMessages[f.step])
}

func (f *reactionsFlow) HandleMessage(session *discordgo.Session, message *discordgo.MessageCreate) Flow {
	if base.IsBotCommand(message, "abort") {
		return NewListeningFlow()
	}

	switch f.step {
	case 0:
		ChannelID := strings.Trim(message.Content, "<#>")
		f.message.ChannelID = ChannelID
		reply := fmt.Sprintf("Alright. The channel is <#%s>", ChannelID)
		base.SendReply(session, reply)
	case 1:
		f.message.UpdateTitleAndDescription(message.Content)
	case 2:
		err := f.message.UpdateColor(message.Content)

		if err != nil {
			reply := fmt.Sprintf("Error converting color: %v\nTry again, or type `none` to leave message without color", err.Error())
			base.SendReply(session, reply)
			return f
		}
		base.SendReply(session, "Perfect. This is how your message will look like:")
		f.postCurrentMessage(session)
	case 3:
		if base.IsBotCommand(message, "done") {
			// TODO: Add option to jump to steps and make corrections
			f.postCurrentMessage(session)
			break
		}

		emoji, r := parseReaction(message.Content)

		// TODO: Make sure role exists

		f.message.Reactions[emoji.Name] = r
		f.message.Emojis[emoji.Name] = emoji
		f.message.MessageID = f.updateCurrentMessage(session)

		err := session.MessageReactionAdd(config.Current.SetupChannelID, f.message.MessageID, emoji.AsReaction())
		if err != nil {
			fmt.Println("Error adding reaction: ", err.Error())
		}

		return f
	case 4:
		if base.IsBotCommand(message, "yes") {
			err := f.postFinalMessage(session)
			if err != nil {
				reply := fmt.Sprintf("Failed to post final message: %v\nType **yes** to try again or **no** to abort and start from the beginning", err.Error())
				base.SendReply(session, reply)
				return f
			} else {
				f.finishConfiguring()
				return NewListeningFlow()

			}
		} else if base.IsBotCommand(message, "no") {
			return NewListeningFlow()
		}
	default:
		fmt.Println("Unknown configuration step: ", f.step)
		f.step = 0
		return f
	}

	f.step += 1
	base.SendReply(session, configurationMessages[f.step])
	return f
}

func (f *reactionsFlow) postCurrentMessage(session *discordgo.Session) {
	err := f.sendMessageToChannel(session, config.Current.SetupChannelID)
	if err != nil {
		fmt.Printf("Failed to post current message: %v", err)
	}
}

func (f *reactionsFlow) sendMessageToChannel(session *discordgo.Session, channelID string) error {
	msg, err := session.ChannelMessageSendComplex(channelID, f.message.ToDiscordMessage())

	if err != nil {
		return err
	}

	for _, emoji := range f.message.Emojis {
		err := session.MessageReactionAdd(msg.ChannelID, msg.ID, emoji.AsReaction())
		if err != nil {
			return err
		}
	}

	f.message.MessageID = msg.ID

	return nil
}

func parseReaction(message string) (emoji base.Emoji, role string) {
	trimmed := strings.Trim(message, " ")
	parts := strings.Fields(trimmed)

	if len(parts) > 2 {
		fmt.Println("Warning: reaction / role message might have too many spaces")
	}

	emoji = base.FindEmojiInMessage(parts[0])
	role = strings.Trim(parts[len(parts)-1], "<@&>")
	return
}

func (f *reactionsFlow) updateCurrentMessage(session *discordgo.Session) string {
	var edit = discordgo.MessageEdit{
		ID:      f.message.MessageID,
		Channel: config.Current.SetupChannelID,
		Embed:   f.message.ToEmbed(),
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

func (f *reactionsFlow) postFinalMessage(session *discordgo.Session) error {
	return f.sendMessageToChannel(session, f.message.ChannelID)
}

func (f *reactionsFlow) finishConfiguring() {
	base.WatchedMessages[f.message.MessageID] = f.message

	err := base.SaveMessages()
	if err != nil {
		fmt.Printf("Failed to save reactions: %v", err.Error())
	}
}
