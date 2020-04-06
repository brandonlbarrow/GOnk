package main

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func streamHandler(s *discordgo.Session, m *discordgo.PresenceUpdate) {
	fmt.Println(m.Game.Details)
}
