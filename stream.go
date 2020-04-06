package main

import (
	"fmt"
	"os"

	"github.com/bwmarrin/discordgo"
)

func streamHandler(s *discordgo.Session, m *discordgo.PresenceUpdate) {

	if m.Presence.Nick != "" {
		if m.Game != nil {
			updateChannel(s, m)
		}
	}
}

func updateChannel(s *discordgo.Session, p *discordgo.PresenceUpdate) {

	if p.Game == nil {
		return
	}

	if p.Game.Name == "Twitch" {
		streamChannel, exists := os.LookupEnv("STREAM_CHANNEL")
		if !exists {
			fmt.Println("Cannot find env variable STREAM_CHANNEL. Please ensure this is set to use streaming alerts.")
			return
		}

		user := getUser(s, p.Presence.User.ID)
		messageBody := formatMessage(user, p.Game.Details, p.Game.URL)
		s.ChannelMessageSend(streamChannel, messageBody)
	}
}

func getUser(s *discordgo.Session, usrID string) string {

	user, err := s.User(usrID)
	if err != nil {
		fmt.Println("Could not find user with id " + usrID)
		os.Exit(1)
	}
	return user.Username
}

func formatMessage(user string, details string, url string) string {

	message := "~STREAM TIME~!\n" + user + " ~went~ ~live~!\n" + details + "\n" + url

	return message
}
