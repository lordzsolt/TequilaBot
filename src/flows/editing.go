package flows

import (
	"fmt"

	"github.com/barkimedes/go-deepcopy"
	"github.com/bwmarrin/discordgo"

	"TequilaBot/src/base"
	"TequilaBot/src/config"
)

func NewEditingFlw() Flow {
	return &editingFlow{}
}

type editingFlow struct {
	step              int
	originalMessageID string
	originalChannelID string
	roleReaction      *base.RoleReaction
}

var editingPrompts = [...]string{
	"Step 1: Which message would you like to edit?",
	"Step 2: To change the title or description, perform it with the original prompt\n" +
		"> `Title | Description`\n" +
		"Or write **skip** to go to the next step",
	"Step 3: Add the emoji, then the name of the role. Type **done** when you are finished.\n" +
		"> Example: `:lemon: @Alcoholic`\n" +
		"Typing the same emoji a second time will override the previous role.\n",
	// "Typing the same emoji a second time will override the previous role.\n" +
	// "To remove an emoji, type **remove** followed by the emoji" +
	// "> Example: `remove :lemon`",
	"Good job! Above you can see how the final message looks like. Type **yes** to post it or **no** to discard it and start again",
}

func (f *editingFlow) Start(session *discordgo.Session) {
	base.PostAllReactions(session)
	base.SendReply(session, editingPrompts[f.step])
}

func (f *editingFlow) HandleMessage(session *discordgo.Session, message *discordgo.MessageCreate) (next Flow) {
	switch f.step {
	case 0:
		err := f.detectReaction(session, message)
		if err != nil {
			fmt.Printf("Failed to start editing flow: %v", err.Error())
		}
		f.nextStep(session)
	case 1:
		if !base.IsBotCommand(message, "skip") {
			f.roleReaction.UpdateTitleAndDescription(message.Content)
			base.EditingExistingRoleReaction(session, *f.roleReaction)
		}

		f.nextStep(session)
	case 2:
		if base.IsBotCommand(message, "done") {
			f.nextStep(session)
		} else {
			emoji := f.roleReaction.UpdateEmojiAndRole(message.Content)
			base.EditingExistingRoleReaction(session, *f.roleReaction)
			err := session.MessageReactionAdd(config.Current.SetupChannelID, f.roleReaction.MessageID, emoji.AsReaction())
			if err != nil {
				fmt.Println("Error adding reaction: ", err.Error())
			}
		}
	case 3:
		f.roleReaction.MessageID = f.originalMessageID
		f.roleReaction.ChannelID = f.originalChannelID
		base.EditingExistingRoleReaction(session, *f.roleReaction)
		for _, emoji := range f.roleReaction.Emojis {
			err := session.MessageReactionAdd(config.Current.SetupChannelID, f.roleReaction.MessageID, emoji.AsReaction())
			if err != nil {
				fmt.Println("Error adding reaction: ", err.Error())
			}
		}
		base.WatchedMessages[f.roleReaction.MessageID] = *f.roleReaction
		err := base.SaveMessages()
		if err != nil {
			fmt.Printf("Failed to save reactions: %v", err.Error())
		}

		return NewListeningFlow()
	default:
		base.SendUnknownMessage(session)
	}

	return f
}

func (f *editingFlow) nextStep(session *discordgo.Session) {
	f.step += 1
	base.SendReply(session, editingPrompts[f.step])
}

func (f *editingFlow) detectReaction(session *discordgo.Session, message *discordgo.MessageCreate) error {
	messageID := message.Message.Content
	originalRoleReaction, ok := base.WatchedMessages[messageID]
	if !ok {
		base.SendReply(session, "Could not find message, please try again or write **abort**")
		return fmt.Errorf("could not find message")
	}
	f.originalMessageID = originalRoleReaction.MessageID
	f.originalChannelID = originalRoleReaction.ChannelID

	copiedRoleReaction, err := deepcopy.Anything(originalRoleReaction)
	if err != nil {
		return err
	}
	roleReaction := copiedRoleReaction.(base.RoleReaction)
	f.roleReaction = &(roleReaction)

	base.SendReply(session, "This is how the message currently looks like:")
	messageID, err = base.SendReactionToChannel(session, roleReaction, config.Current.SetupChannelID, true)

	if err != nil {
		return err
	}

	f.roleReaction.MessageID = messageID
	return nil
}
