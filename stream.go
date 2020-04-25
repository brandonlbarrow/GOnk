package main

import (
	"fmt"
	"os"

	"github.com/bwmarrin/discordgo"
)

func streamHandler(s *discordgo.Session, m *discordgo.PresenceUpdate) {

	guildID, exists := os.LookupEnv("GUILD_ID")
	if !exists {
		fmt.Println("Cannot find env variable GUILD_ID. Please ensure this is set to use gonk.")
		os.Exit(1)
	}

	if m.Game == nil {
		return
	}

	if m.GuildID == guildID {
		if m.Presence.Game.Type == 1 {
			updateChannel(s, m)
		}
	}
}

func updateChannel(s *discordgo.Session, p *discordgo.PresenceUpdate) {

	if p.Game.Name == "Twitch" {
		streamChannel, exists := os.LookupEnv("STREAM_CHANNEL")
		if !exists {
			fmt.Println("Cannot find env variable STREAM_CHANNEL. Please ensure this is set to use streaming alerts.")
			return
		}

		user := getUser(s, p.Presence.User.ID)
		if p.Nick != "" {
			user = p.Nick
		}

		messageBody := formatMessage(user, p.Game.State, p.Game.Details, p.Game.URL)
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

func formatMessage(user string, assets string, details string, url string) string {

	message := "~STREAM TIME~!\n" + "**" + user + "**" + " ~went live with~ " + "**" + assets + "**" + "!\n" + details + "\n" + url

	return message
}
