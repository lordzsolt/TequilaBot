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

var reactionsPrompts = [...]string{
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
	step            int
	targetChannelID string
	message         base.RoleReaction
}

func (f *reactionsFlow) Start(session *discordgo.Session) {
	f.message = base.NewRoleReaction()
	f.message.ChannelID = config.Current.SetupChannelID
	base.SendReply(session, reactionsPrompts[f.step])
}

func (f *reactionsFlow) HandleMessage(session *discordgo.Session, message *discordgo.MessageCreate) Flow {
	if base.IsBotCommand(message, "abort") {
		return NewListeningFlow()
	}

	switch f.step {
	case 0:
		f.targetChannelID = strings.Trim(message.Content, "<#>")
		reply := fmt.Sprintf("Alright. The channel is <#%s>", f.targetChannelID)
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

		emoji := f.message.UpdateEmojiAndRole(message.Content)
		base.EditingExistingRoleReaction(session, f.message)
		err := session.MessageReactionAdd(config.Current.SetupChannelID, f.message.MessageID, emoji.AsReaction())
		if err != nil {
			fmt.Println("Error adding reaction: ", err.Error())
		}

		return f
	case 4:
		if base.IsBotCommand(message, "no") {
			return NewListeningFlow()
		}
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
		}
		return f
	default:
		fmt.Println("Unknown configuration step: ", f.step)
		f.step = 0
		return f
	}

	f.step += 1
	base.SendReply(session, reactionsPrompts[f.step])
	return f
}

func (f *reactionsFlow) postCurrentMessage(session *discordgo.Session) {
	messageID, err := base.SendReactionToChannel(session, f.message, config.Current.SetupChannelID, true)
	if err != nil {
		fmt.Printf("Failed to post current message: %v", err)
	}
	f.message.MessageID = messageID
}

func (f *reactionsFlow) postFinalMessage(session *discordgo.Session) error {
	messageID, err := base.SendReactionToChannel(session, f.message, f.targetChannelID, true)
	if err != nil {
		fmt.Printf("Failed to post current message: %v", err)
	} else {
		f.message.ChannelID = f.targetChannelID
		f.message.MessageID = messageID
	}
	return err
}

func (f *reactionsFlow) finishConfiguring() {
	base.WatchedMessages[f.message.MessageID] = f.message

	err := base.SaveMessages()
	if err != nil {
		fmt.Printf("Failed to save reactions: %v", err.Error())
	}
}
