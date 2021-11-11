package info

import (
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
)

const (
	infoCmd    = "!info"
	gonkRepo   = "https://github.com/brandonlbarrow/GOnk"
	maintainer = "brandonlbarrow"
	details    = `
	GOnk is a Discord bot that announces Twitch streams and tells you how to make alcoholic drinks, among other things. GOnk is maintained by Brandoid and is written in the Go programming language. ~GoNk~
	`
)

// Handler is an object representing the application's runtime state
type Handler struct {
	version    string
	repo       string
	maintainer string
	start      time.Time
	details    string
}

type HandlerOption func(h *Handler)

func WithVersion(version string) HandlerOption {
	return func(h *Handler) {
		h.version = version
	}
}

func NewHandler(opts ...HandlerOption) *Handler {
	h := &Handler{
		version:    "1.0.0",
		repo:       gonkRepo,
		maintainer: maintainer,
		start:      time.Now(),
		details:    details,
	}
	for _, opt := range opts {
		opt(h)
	}
	return h
}

func (h *Handler) Handle(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}
	if strings.HasPrefix(m.Content, infoCmd) {
		s.ChannelMessageSendEmbed(m.ChannelID, h.infoEmbed())
	}
}

func (h *Handler) infoEmbed() *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title:     "GOnk Info",
		Timestamp: time.Now().Format(time.RFC3339),
		Color:     0x33ff33,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Version",
				Value:  h.version,
				Inline: true,
			},
			{
				Name:   "Uptime",
				Value:  time.Since(h.start).String(),
				Inline: true,
			}, {
				Name:   "GitHub",
				Value:  h.repo,
				Inline: true,
			}, {
				Name:   "CPU Util",
				Value:  getCPUPercent(),
				Inline: true,
			}, {
				Name:   "Memory Util",
				Value:  getMemoryUsage(),
				Inline: true,
			},
			{
				Name:   "About Me",
				Value:  h.details,
				Inline: false,
			},
		},
	}
}

func getCPUPercent() string {
	val, err := cpu.Percent(0, false)
	if err != nil {
		return fmt.Sprintf("%v", err)
	}
	if len(val) == 0 {
		return ""
	}
	return fmt.Sprintf("%.2f%%", val[0])
}

func getMemoryUsage() string {
	val, err := mem.VirtualMemory()
	if err != nil {
		return fmt.Sprintf("%v", err)
	}
	if val == nil {
		return ""
	}
	return fmt.Sprintf("%.2f%%", val.UsedPercent)
}
