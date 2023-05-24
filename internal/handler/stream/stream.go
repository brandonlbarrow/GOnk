package stream

import (
	"sync"

	"github.com/brandonlbarrow/gonk/v2/internal/db"
	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
)

const (
	statusStreamingKey = "streaming"
)

// Handler receives PresenceUpdate events from the Discord API and handles them. Streaming notifications will be sent to the channelID and guildID for any user.
// If a userID is supplied, then stream notification events will only be sent for events matching the userID.
type Handler struct {
	logger *logrus.Logger
	repo   db.ServerRepo
}

// Sessioner is used by *discordgo.Session objects
type Sessioner interface {
	ChannelMessageSend(channelID string, content string) (*discordgo.Message, error)
	User(userID string) (st *discordgo.User, err error)
}

// NewHandler creates an instance of *Handler.
func NewHandler(opts ...HandlerOption) *Handler {
	logger := logrus.New()
	h := &Handler{
		logger: logger,
	}
	for _, opt := range opts {
		opt(h)
	}
	return h
}

// HandlerOption is an option passed to NewHandler to set struct fields.
type HandlerOption func(m *Handler)

// WithLogger sets the logrus instance for the Handler.
func WithLogger(logger *logrus.Logger) HandlerOption {
	return func(m *Handler) {
		m.logger = logger
	}
}

func WithRepo(repo db.ServerRepo) HandlerOption {
	return func(m *Handler) {
		m.repo = repo
	}
}

// Handle implements the Discordgo API and receives PresenceUpdate events. See the documentation for *discordgo.PresenceUpdate to see all available fields. Handle will parse this event and send a notification that a user started streaming
// on Twitch if their Discord user account has a Twitch notification set up, they are set to show their Game (currently playing) status on Discord, and they go from a non-streaming state to a streaming state.
// "Game" is Discord's term for any game or streaming session the user may be playing/hosting. Game is usually a video game or some program that Discord assumes is a game, and Game becomes a streaming status when the user goes live and has
// a streaming platform such as Twitch integrated with their account.

func (h *Handler) Handle(s *discordgo.Session, m *discordgo.MessageCreate) {
	h.streamHandler(s, m)
}

// streamHandler configures Gonk's stream handler for the Discord server and stores the configuration in the repo
func (m *Handler) streamHandler(s Sessioner, p *discordgo.MessageCreate) {

	guildID := extractGuildIDFromMessageCreate(p)
	serverConfig, err := m.repo.GetServerByID(guildID)
	if err != nil {
		m.logger.Errorf("error getting server config: %w", err)
		return
	} else if serverConfig != nil {

	}
}

func (m *Handler) getUser(s Sessioner, usrID string) (string, error) {

	user, err := s.User(usrID)
	if err != nil {
		return "", err
	}
	return user.Username, nil
}

func formatMessage(user string, assets string, details string, url string) string {

	message := "~STREAM TIME~!\n" + "**" + user + "**" + " ~went live with~ " + "**" + assets + "**" + "!\n" + details + "\n" + url

	return message
}

func extractGuildIDFromMessageCreate(m *discordgo.MessageCreate) string {
	if m != nil && m.GuildID != "" {
		return m.GuildID
	}
	if m.Member != nil && m.Member.GuildID != "" {
		return m.Member.GuildID
	}
	if m.MessageReference != nil && m.MessageReference.GuildID != "" {
		return m.MessageReference.GuildID
	}
}
