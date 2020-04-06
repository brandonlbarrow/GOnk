package main

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func chatHandler(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Content == "whoami" {
		user, err := s.User(m.Author.ID)
		if err != nil {
			fmt.Println(err)
		}
		s.ChannelMessageSend(m.ChannelID, user.Mention())
	}
}

func listMembers(s *discordgo.Session, g string) []*discordgo.Member {

	memberList, err := s.GuildMembers(g, "0", 100)
	if err != nil {
		fmt.Println(err)
	}

	return memberList
}
