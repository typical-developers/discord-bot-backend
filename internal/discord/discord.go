package discord

import (
	"os"

	"github.com/bwmarrin/discordgo"
)

var (
	Client *discordgo.Session
)

func init() {
	c, err := discordgo.New("Bot " + os.Getenv("DISCORD_TOKEN"))
	if err != nil {
		panic(err)
	}

	Client = c

	Client.StateEnabled = true
}
