package handler

import "github.com/bwmarrin/discordgo"

type PresenceUpdateHandler interface {
	Handle(s *discordgo.Session, p *discordgo.PresenceUpdate) error
}

type MessageCreateHandler interface {
	Handle(s *discordgo.Session, m *discordgo.MessageCreate) error
}
