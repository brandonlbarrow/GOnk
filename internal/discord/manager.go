package discord

import (
	"context"
	"fmt"

	"github.com/bwmarrin/discordgo"
)

// Manager manages the discord session for a server (guild).
type Manager struct {
	discordSession *discordgo.Session
	guildID        string
}

// ManagerOption is passed as a slice to NewManager to set the configuration of the Manager.
type ManagerOption func(c *Manager)

func NewManager(opts ...ManagerOption) *Manager {
	m := &Manager{}
	for _, opt := range opts {
		opt(m)
	}
	return m
}

// WithGuildID sets the guildID for the *Manager.
func WithGuildID(guildID string) ManagerOption {
	return func(m *Manager) {
		m.guildID = guildID
	}
}

// SessionArgs are the arguments to be assigned to the *discordgo.Session.
type SessionArgs struct {
	StateEnabled bool
	LogLevel     int
	Intents      *discordgo.Intent
}

// Returns a new SessionArgs object using defaults.
// https://discord.com/developers/docs/topics/gateway#gateway-intents
func NewSessionArgsWithDefaults() SessionArgs {
	return SessionArgs{
		LogLevel:     discordgo.LogError,
		StateEnabled: true,
		Intents:      discordgo.MakeIntent(discordgo.IntentsGuildPresences | discordgo.IntentsGuildMessages | discordgo.IntentsGuildMessageReactions),
	}
}

// MustWithSession takes a token string and returns a fully initialized Discord *Session object. If the session cannot be established, this function will panic with the associated error.
func MustWithSession(token string, args SessionArgs) ManagerOption {
	return func(m *Manager) {
		session, err := discordgo.New("Bot " + token)
		if err != nil {
			panic(err)
		}
		session.LogLevel = args.LogLevel
		session.StateEnabled = args.StateEnabled
		session.Identify.Intents = args.Intents
		m.discordSession = session
	}
}

// AddHandler adds supplied Discord session handler functions to the session.
func (m *Manager) AddHandler(handler interface{}) {
	m.discordSession.AddHandler(handler)
}

// Run opens a websocket connection to Discord's API and errors if it runs into a problem doing so.
func (m *Manager) Run(ctx context.Context) error {
	err := m.discordSession.Open()
	if err != nil {
		return fmt.Errorf("error opening discord session: %w", err)
	}
	return nil
}
