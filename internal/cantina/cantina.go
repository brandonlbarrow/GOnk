package cantina

import (
	"log"
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// Handler ...
func Handler(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	if strings.Contains(m.Content, os.Getenv("CANTINA_LISTEN_TEXT")) {
			_, err := s.ChannelMessageSend(m.ChannelID, os.Getenv("CANTINA_URL"))
		if err != nil {
			log.Fatalf("failed sending message to %s with content %s", m.ChannelID, os.Getenv("CANTINA_URL"))
		}
	}



}