package main

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func chatHandler(s *discordgo.Session, m *discordgo.MessageCreate) {

	GuildID := getGuildID(s)

	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Content == "whoami" {
		usr := m.Author.ID
		usr2, err := s.User(usr)
		if err != nil {
			fmt.Println(err)
		}
		s.ChannelMessageSend(m.ChannelID, usr2.Mention())
	}

	if m.Content == "guildid" {
		s.ChannelMessageSend(m.ChannelID, GuildID)
	}

	if m.Content == "start" {
		s.UpdateStreamingStatus(0, "Twitch", "http://www.twitch.tv/gonk")
	}

	if m.Content == "end" {
		s.UpdateStreamingStatus(0, "", "")
	}

	if m.Content == "list" {
		members := listMembers(s, GuildID)
		for _, member := range members {
			s.ChannelMessageSend(m.ChannelID, member.User.Username)
		}
	}
}

func listMembers(s *discordgo.Session, g string) []*discordgo.Member {

	memberList, err := s.GuildMembers(g, "0", 100)
	if err != nil {
		fmt.Println(err)
	}

	return memberList
}
