package bot

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

var (
	config                 *Config
	currentState           state
	configurationStep      int
	messageBeingConfigured *roleReactionMessage
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

func Start(c *Config) error {
	config = c

	session, err := discordgo.New("Bot " + config.Token)
	if err != nil {
		return err
	}

	session.AddHandler(messageCreated)
	session.AddHandler(reactionAdd)

	session.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentsGuildMessageReactions

	err = session.Open()
	if err != nil {
		return err
	}

	fmt.Println("Bot is running")

	test(session)

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)

	<-sc

	return session.Close()
}

func test(session *discordgo.Session) {
	// startConfiguring(session)
	// messageBeingConfigured.title = "Lemon is gay"
	// messageBeingConfigured.description = "Uhum..."
	// messageBeingConfigured.messageID = "834914029486604298"
	// configurationStep = 3
}

func isBotCommand(message *discordgo.MessageCreate, command string) bool {
	return strings.ToLower(message.Content) == command
}

func sendReply(session *discordgo.Session, message string) {
	_, err := session.ChannelMessageSend(config.SetupChannelID, message)

	if err != nil {
		fmt.Println("Failed to send message: ", err.Error())
	}
}
