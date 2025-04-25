package discord

import (
	"os"

	"github.com/bwmarrin/discordgo"
)

type DiscordClient struct {
	sess *discordgo.Session
}

func NewDiscordClient() *DiscordClient {
	s, err := discordgo.New("Bot " + os.Getenv("BOT_TOKEN"))
	if err != nil {
		panic("Failed to create a new DiscordClient")
	}
	return &DiscordClient{
		sess: s,
	}
}
