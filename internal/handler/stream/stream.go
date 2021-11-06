package stream

import (
	"fmt"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
)

// StreamList map of maps containing all stream activity
var StreamList = make(map[string]map[string]bool)

// Handler receives PresenceUpdate events from the Discord API and handles them
type Handler struct {
	logger *logrus.Logger
	// channelID is the channel that stream events will be sent to.
	channelID string
	guildID   string
}

// NewHandler creates an instance of *Handler.
func NewHandler(channelID, guildID string, logger *logrus.Logger) *Handler {
	return &Handler{logger: logger, channelID: channelID, guildID: guildID}
}

func (m *Handler) Handle(s *discordgo.Session, p *discordgo.PresenceUpdate) {

	m.streamHandler(s, p)
	message := ""
	msg, err := s.ChannelMessageSend(m.channelID, message)
	if err != nil {
		m.logger.Debugf("error sending message.\nPresence: %v\nError: %w\nChannelID: %s\nGuildID: %s", p, err, m.channelID, m.guildID)
		return
	}
	m.logger.Debugf("message successfully sent.\nMessage: %v\nChannelID: %s\nGuildID: %s", msg, m.channelID, m.guildID)
	return
}
func (m *Handler) streamHandler(s *discordgo.Session, p *discordgo.PresenceUpdate) {
	m.logger.Infof("invoking stream handler.\n session: %s\n presenceUpdate: %s", s, p)

	guildID, exists := os.LookupEnv("GUILD_ID")
	if !exists {
		m.logger.Error("Cannot find env variable GUILD_ID. Please ensure this is set to use gonk.")
		return
	}

	streamChannel, exists := os.LookupEnv("STREAM_CHANNEL")
	if !exists {
		logrus.Error("Cannot find env variable STREAM_CHANNEL. Please ensure this is set to use streaming alerts.")
		return
	}

	if !validateGuildID(p, guildID) {
		logrus.Errorf("cannot validate guild id: %s", guildID)
		return
	}

	userID := p.Presence.User.ID
	_, ok := StreamList[userID]
	if !ok {
		StreamList[userID] = map[string]bool{"streaming": false}
	}

	if p.Game == nil {
		StreamList[userID]["streaming"] = false
		fmt.Println(StreamList[userID])
		fmt.Println("No game")
		return
	}

	if p.Game.Type == 1 {
		if StreamList[userID]["streaming"] {
			fmt.Println("already streaming")
			return
		} else {
			StreamList[userID]["streaming"] = true
			fmt.Println(StreamList[userID])
			fmt.Println("Stream started")
		}
	}

	if p.Game.Type != 1 {
		if StreamList[userID]["streaming"] == false {
			fmt.Println("already not streaming")
		} else {
			StreamList[userID]["streaming"] = false
			fmt.Println(StreamList[userID])
			fmt.Println("Stream ended or not streaming")
		}
	}

	if StreamList[userID]["streaming"] == true {
		user := getUser(s, p.Presence.User.ID)
		if p.Nick != "" {
			user = p.Nick
		}
		messageBody := formatMessage(user, p.Game.State, p.Game.Details, p.Game.URL)
		s.ChannelMessageSend(streamChannel, messageBody)
	}
}

func getUser(s *discordgo.Session, usrID string) string {

	user, err := s.User(usrID)
	if err != nil {
		fmt.Println("Could not find user with id " + usrID)
		os.Exit(1)
	}

	return user.Username
}

func formatMessage(user string, assets string, details string, url string) string {

	message := "~STREAM TIME~!\n" + "**" + user + "**" + " ~went live with~ " + "**" + assets + "**" + "!\n" + details + "\n" + url

	return message
}

func validateGuildID(p *discordgo.PresenceUpdate, g string) bool {

	if p.GuildID != g {
		return false
	}
	return true
}
