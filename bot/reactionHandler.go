package bot

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func reactionAdd(session *discordgo.Session, reaction *discordgo.MessageReactionAdd) {
	fmt.Println(reaction)
}
