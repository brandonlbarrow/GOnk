package cantina

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

// Handler ...
func Handler(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	if strings.Contains(m.Content, os.Getenv("CANTINA_LISTEN_TEXT")) {
		url := fmt.Sprintf("%s?seed=%s", os.Getenv("CANTINA_URL"), strconv.FormatInt((time.Now().UnixNano()), 10))
		fmt.Println(url)
		fmt.Println(m.ChannelID)
		_, err := s.ChannelMessageSend(m.ChannelID, url)
		if err != nil {
			log.Fatalf("failed sending message to %s with content %s", m.ChannelID, os.Getenv("CANTINA_URL"))
		}
	}
}
