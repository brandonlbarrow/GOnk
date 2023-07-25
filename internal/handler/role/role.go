package role

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

const (
	ffxivEmoji    = "kupo"
	ffxivRoleID   = "849859950637350922"
	padawanEmoji  = "jediface"
	padawanRoleID = "573973415480918216"
)

type Handler struct {
	guildID string
}

type HandlerOption func(h *Handler)

func WithGuildID(guildID string) HandlerOption {
	return func(h *Handler) {
		h.guildID = guildID
	}
}

func NewHandler(opts ...HandlerOption) *Handler {
	h := &Handler{}
	for _, opt := range opts {
		opt(h)
	}
	return h
}

func (h *Handler) Handle(s *discordgo.Session, m *discordgo.MessageReactionAdd) {
	if m.Emoji.Name == ffxivEmoji {
		if err := s.GuildMemberRoleAdd(m.GuildID, m.UserID, ffxivRoleID); err != nil {
			fmt.Println(err)
		}
		if err := s.MessageReactionRemove(m.ChannelID, m.MessageID, m.Emoji.APIName(), m.UserID); err != nil {
			fmt.Println(err)
		}
		return
	} else if m.Emoji.Name == padawanEmoji {
		if err := s.GuildMemberRoleAdd(m.GuildID, m.UserID, padawanRoleID); err != nil {
			fmt.Println(err)
		}
		if err := s.MessageReactionRemove(m.ChannelID, m.MessageID, m.Emoji.APIName(), m.UserID); err != nil {
			fmt.Println(err)
		}
	} else {
		if err := s.MessageReactionRemove(m.ChannelID, m.MessageID, m.Emoji.APIName(), m.UserID); err != nil {
			fmt.Println(err)
		}
	}
}
