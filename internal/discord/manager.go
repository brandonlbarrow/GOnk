package discord

import (
	"context"
	"errors"
	"fmt"

	"github.com/bwmarrin/discordgo"
)

// Manager manages the discord session.
type Manager struct {
	discordSession *discordgo.Session
	discordToken   string
}

type ManagerOption func(c *Manager)

func NewManager(opts ...ManagerOption) *Manager {
	m := &Manager{}
	for _, opt := range opts {
		opt(m)
	}
	return m
}

func WithToken(token string) ManagerOption {
	return func(m *Manager) {
		m.discordToken = token
	}
}

func (m *Manager) NewDiscordSession() error {

	session, err := discordgo.New("Bot " + m.discordToken)
	if err != nil {
		return errors.New("error initializing discord session")
	}

	session.StateEnabled = true
	// https://discord.com/developers/docs/topics/gateway#gateway-intents
	session.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsGuildPresences | discordgo.IntentsGuildMessages | discordgo.IntentsGuildMessageReactions)

	m.discordSession = session
	return nil
}

func (m *Manager) AddHandler(handlers ...interface{}) {
	for _, handler := range handlers {
		m.discordSession.AddHandler(handler)
	}
}

func (m *Manager) Run(ctx context.Context) error {
	err := m.discordSession.Open()
	if err != nil {
		return fmt.Errorf("Error opening discord session: %w", err)
	}
	return nil
}

func (m *Manager) Token() string {
	return m.discordToken
}
