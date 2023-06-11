package stream

import (
	"sync"

	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
)

const (
	statusStreamingKey = "streaming"
)

// Handler receives PresenceUpdate events from the Discord API and handles them. Streaming notifications will be sent to the channelID and guildID for any user.
// If a userID is supplied, then stream notification events will only be sent for events matching the userID.
type Handler struct {
	logger      *logrus.Logger
	guildID     string
	channelID   string
	userID      string
	streamerMap *streamerMap
}

type streamerMap struct {
	streamList map[string]streamerStatus
	lock       sync.RWMutex
}

func newStreamerMap() *streamerMap {
	streamList := make(map[string]streamerStatus)
	return &streamerMap{
		lock:       sync.RWMutex{},
		streamList: streamList,
	}
}

type streamerStatus map[string]bool

func (s *streamerMap) getStreamList() map[string]streamerStatus {
	return s.streamList
}

func (s *streamerMap) userIsStreaming(userID string) bool {
	s.lock.Lock()
	defer s.lock.Unlock()
	if _, ok := s.streamList[userID]; ok {
		return s.streamList[userID][statusStreamingKey]
	}
	return false
}

func (s *streamerMap) setUserStreamStatus(userID string, streaming bool) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.streamList[userID][statusStreamingKey] = streaming
}

// Sessioner is used by *discordgo.Session objects
type Sessioner interface {
	ChannelMessageSend(string, string, ...discordgo.RequestOption) (*discordgo.Message, error)
	User(string, ...discordgo.RequestOption) (*discordgo.User, error)
}

// NewHandler creates an instance of *Handler.
func NewHandler(opts ...HandlerOption) *Handler {
	streamerMap := newStreamerMap()
	logger := logrus.New()
	h := &Handler{
		logger:      logger,
		streamerMap: streamerMap,
	}
	for _, opt := range opts {
		opt(h)
	}
	return h
}

// HandlerOption is an option passed to NewHandler to set struct fields.
type HandlerOption func(m *Handler)

// WithGuildID sets the discord guild ID for the Handler.
func WithGuildID(guildID string) HandlerOption {
	return func(m *Handler) {
		m.guildID = guildID
	}
}

// WithChannelID sets the discord channel ID for the Handler to send streaming events to.
func WithChannelID(channelID string) HandlerOption {
	return func(m *Handler) {
		m.channelID = channelID
	}
}

// WithLogger sets the logrus instance for the Handler.
func WithLogger(logger *logrus.Logger) HandlerOption {
	return func(m *Handler) {
		m.logger = logger
	}
}

// WithUserID sets the discord user ID for the Handler to send streaming events for. If this is set, only the user ID provided will trigger notifications when they go live.
func WithUserID(userID string) HandlerOption {
	return func(m *Handler) {
		m.userID = userID
	}
}

// Handle implements the Discordgo API and receives PresenceUpdate events. See the documentation for *discordgo.PresenceUpdate to see all available fields. Handle will parse this event and send a notification that a user started streaming
// on Twitch if their Discord user account has a Twitch notification set up, they are set to show their Game (currently playing) status on Discord, and they go from a non-streaming state to a streaming state.
// "Game" is Discord's term for any game or streaming session the user may be playing/hosting. Game is usually a video game or some program that Discord assumes is a game, and Game becomes a streaming status when the user goes live and has
// a streaming platform such as Twitch integrated with their account.
func (m *Handler) Handle(s *discordgo.Session, p *discordgo.PresenceUpdate) {

	m.logger.WithFields(presenceUpdateFields(p)).Info("presenceUpdate user info:")
	m.logger.WithFields(logrus.Fields{"streamList": m.streamerMap.getStreamList()}).Info("current stream list:")
	m.streamHandler(s, p)
}

// streamHandler references an in-memory map in Handler to keep track of user streaming state. This map is populated as events come in.
func (m *Handler) streamHandler(s Sessioner, p *discordgo.PresenceUpdate) {

	// if the PresenceUpdate object is empty or the User is nil, return false. We don't care about events that don't have these.
	if !validatePresenceUpdateObject(p) {
		m.logger.WithFields(logrus.Fields{"presenceUpdateObject": p}).Debug("presenceUpdate failed validation, skipping")
		return
	}

	// validate that the PresenceUpdate's server id matches the one Gonk has been configured to operate in, otherwise skip it
	if !validateGuildID(p, m.guildID) {
		m.logger.WithFields(logrus.Fields{"providedGuildID": m.guildID, "eventGuildID": p.GuildID}).Debug("guild ID does not match PresenceUpdate, skipping")
		return
	}

	// validation for single user mode
	if !validateUserID(p, m.userID) {
		m.logger.WithFields(logrus.Fields{"providedUserID": m.userID, "eventUser": p.User}).Debug("user ID does not match PresenceUpdate, skipping")
		return
	}

	// get the userID and initialize their streaming state as false if they don't already exist in the map of streamerIDs to streamerStatus
	userID := p.Presence.User.ID
	_, ok := m.streamerMap.streamList[userID]
	if !ok {
		m.streamerMap.streamList[userID] = map[string]bool{"streaming": false}
	}

	if len(p.Presence.Activities) == 0 {
		m.streamerMap.setUserStreamStatus(userID, false)
	}
	for i, activity := range p.Presence.Activities {
		m.logger.WithFields(logrus.Fields{"index": i, "activityName": activity.Name, "activityType": activity.Type}).Info("activities")
		if activity.Type == discordgo.ActivityTypeStreaming {
			streaming := m.streamerMap.userIsStreaming(userID)
			if streaming {
				m.logger.WithFields(logrus.Fields{"userID": userID, "gameType": activity.Type, "streamingStatus": streaming}).Debug("no change")
			} else {
				m.streamerMap.setUserStreamStatus(userID, true)
				m.logger.WithFields(logrus.Fields{"userID": userID, "gameType": activity.Type, "streamingStatus": m.streamerMap.userIsStreaming(userID)}).Info("user has started streaming.")
			}
		}

		if activity.Type != discordgo.ActivityTypeStreaming {
			streaming := m.streamerMap.userIsStreaming(userID)
			if !streaming {
				m.logger.WithFields(logrus.Fields{"userID": userID, "gameType": activity.Type, "streamingStatus": streaming}).Debug("no change")
			} else {
				m.streamerMap.setUserStreamStatus(userID, false)
				m.logger.WithFields(logrus.Fields{"userID": userID, "gameType": activity.Type, "streamingStatus": m.streamerMap.userIsStreaming(userID)}).Info("Stream ended, or not streaming anymore.")
			}
		}
		if m.streamerMap.userIsStreaming(userID) {
			user, err := m.getUser(s, p.Presence.User.ID)
			if err != nil {
				m.logger.WithFields(logrus.Fields{"userID": userID}).Error("could not find username from supplied user id, cannot send streaming message.")
				return
			}
			messageBody := formatMessage(user, activity.State, activity.Details, activity.URL)
			s.ChannelMessageSend(m.channelID, messageBody)
		}
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

func validateGuildID(p *discordgo.PresenceUpdate, g string) bool {
	return p.GuildID == g
}

func validateUserID(p *discordgo.PresenceUpdate, u string) bool {
	if u == "" {
		return true
	}
	if p.User != nil {
		return p.User.ID == u
	}
	return false
}

func presenceUpdateFields(p *discordgo.PresenceUpdate) logrus.Fields {
	if p == nil {
		return logrus.Fields{}
	}
	baseFields := logrus.Fields{
		"guildID": p.GuildID,
		"status":  p.Status,
	}
	if p.User != nil {
		userFields := logrus.Fields{
			"username": p.User.Username,
			"id":       p.User.ID,
		}
		for k, v := range userFields {
			baseFields[k] = v
		}
	}
	return baseFields
}

func validatePresenceUpdateObject(p *discordgo.PresenceUpdate) bool {
	if p == nil {
		return false
	}
	if p.Presence.User == nil {
		return false
	}
	return true
}
