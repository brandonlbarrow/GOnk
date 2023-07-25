package discord

import (
	"context"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// DiscordgoLogLevel is a wrapper type for the discordgo LogLevel constants.
type DiscordgoLogLevel int

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
	Intents      *discordgo.Intent // https://discord.com/developers/docs/topics/gateway#gateway-intents

}

type SessionArg func(s *SessionArgs)

// Returns a new SessionArgs object using defaults. A slice of SessionArg can be supplied to override the defaults.
func NewSessionArgs(args ...SessionArg) *SessionArgs {
	s := &SessionArgs{
		LogLevel:     discordgo.LogError,
		StateEnabled: true,
		//	Intents:      discordgo.MakeIntent(discordgo.IntentsGuildPresences | discordgo.IntentsGuildMessages | discordgo.IntentsGuildMessageReactions | discordgo.IntentsDirectMessageReactions),
	}
	for _, arg := range args {
		arg(s)
	}
	return s
}

// WithLogLevel sets the DiscordgoLogLevel for the *SessionArgs.
func WithLogLevel(logLevel DiscordgoLogLevel) SessionArg {
	return func(s *SessionArgs) {
		s.LogLevel = int(logLevel)
	}
}

// MustWithSession takes a token string and returns a fully initialized Discord *Session object. If the session cannot be established, this function will panic with the associated error.
func MustWithSession(token string, args *SessionArgs) ManagerOption {
	return func(m *Manager) {
		session, err := discordgo.New("Bot " + token)
		if err != nil {
			panic(err)
		}
		session.LogLevel = args.LogLevel
		session.StateEnabled = args.StateEnabled
		//		session.Identify.Intents = args.Intents
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

// DiscordLogLevelFromString takes a string representation of a logging level and returns the appropriate DiscordgoLogLevel.
func DiscordLogLevelFromString(level string) DiscordgoLogLevel {
	var (
		LevelDebug   DiscordgoLogLevel = DiscordgoLogLevel(discordgo.LogDebug)
		LevelInfo    DiscordgoLogLevel = DiscordgoLogLevel(discordgo.LogInformational)
		LevelError   DiscordgoLogLevel = DiscordgoLogLevel(discordgo.LogError)
		LevelWarning DiscordgoLogLevel = DiscordgoLogLevel(discordgo.LogWarning)
	)
	switch strings.ToLower(level) {
	case "debug":
		return LevelDebug
	case "info":
		return LevelInfo
	case "error":
		return LevelError
	case "warning":
		return LevelWarning
	default:
		return LevelError
	}
}
