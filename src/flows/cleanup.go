package flows

import (
	"fmt"

	"github.com/bwmarrin/discordgo"

	"TequilaBot/src/base"
	"TequilaBot/src/config"
)

func NewCleanupFlow() Flow {
	return &cleanupFlow{}
}

type cleanupFlow struct {
	step                 int
	roleReactionToDelete *base.RoleReaction
}

func (f *cleanupFlow) Start(session *discordgo.Session) {
	f.CleanupUnknownMessages(session)
	base.PostAllReactions(session)
	base.SendReply(session, "Which message ID would you like to delete? Please input the message ID, or **abort**")

	f.step = 0
	f.roleReactionToDelete = nil
}

func (f *cleanupFlow) HandleMessage(session *discordgo.Session, message *discordgo.MessageCreate) (next Flow) {
	switch f.step {
	case 0:
		err := f.detectReactionAndConfirmDeletion(session, message)
		if err != nil {
			fmt.Printf("Failed to extract message ID from %v, error %v", message.Content, err.Error())
		} else {
			f.step = 1
		}
	case 1:
		if base.IsBotCommand(message, "yes") {
			f.deleteMessage(session)
		}

		base.SendReply(session, "Which message ID would you like to delete? Please input the message ID, or **abort**")
		f.step = 0
		f.roleReactionToDelete = nil
	}

	return f
}

func (f *cleanupFlow) CleanupUnknownMessages(session *discordgo.Session) {
	base.SendReply(session, "Cleaning up unknown messages")

	var messagesToDelete []base.RoleReaction
	for _, message := range base.WatchedMessages {
		_, err := session.ChannelMessage(message.ChannelID, message.MessageID)
		if err != nil {
			messagesToDelete = append(messagesToDelete, message)
			errorMessage := fmt.Errorf("Failed to get message with ID %v: %v\n", message.MessageID, err.Error())
			fmt.Print(errorMessage)
			base.SendReply(session, errorMessage.Error()+"; Deleting")
		}
	}

	if len(messagesToDelete) == 0 {
		base.SendReply(session, "Cleanup finished, no unknown messages found")
	} else {
		for _, messagesToDelete := range messagesToDelete {
			delete(base.WatchedMessages, messagesToDelete.MessageID)
		}

		err := base.SaveMessages()
		if err != nil {
			fmt.Printf("Failed to save reactions: %v", err.Error())
		}

		base.SendReply(session, "Cleanup finished.")
	}
}

func (f *cleanupFlow) detectReactionAndConfirmDeletion(session *discordgo.Session, message *discordgo.MessageCreate) error {
	messageID := message.Message.Content
	roleReaction, ok := base.WatchedMessages[messageID]
	if !ok {
		base.SendReply(session, "Could not find message, please try again or write **abort**")
		return fmt.Errorf("could not find message")
	}

	_, err := base.SendReactionToChannel(session, roleReaction, config.Current.SetupChannelID, false)
	if err != nil {
		return err
	}
	base.SendReply(session, "Are you sure you want to delete this message? This cannot be undone. Type **yes** or **no**")
	f.roleReactionToDelete = &roleReaction
	return nil
}

func (f *cleanupFlow) deleteMessage(session *discordgo.Session) {
	err := session.ChannelMessageDelete(f.roleReactionToDelete.ChannelID, f.roleReactionToDelete.MessageID)
	if err != nil {
		errorMessage := fmt.Sprintf("Failed to delete message: %v", err.Error())
		base.SendReply(session, errorMessage)
	} else {
		base.SendReply(session, "Successfully deleted")
	}

	delete(base.WatchedMessages, f.roleReactionToDelete.MessageID)
	err = base.SaveMessages()
	if err != nil {
		fmt.Printf("Failed to save reactions: %v", err.Error())
	}
}
