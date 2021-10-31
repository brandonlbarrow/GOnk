package stream

import (
	"fmt"
	"os"

	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

// StreamList map of maps containing all stream activity
var StreamList = make(map[string]map[string]bool)

type StreamManager struct {
	PresenceHandler func(*discordgo.Session, *discordgo.PresenceUpdate)
	StreamStateMap  map[string]bool
}

var logger log.Logger

func (s *StreamManager) shiftStreamState(userID string, streamPresence int) {
	switch s.StreamStateMap[userID] {
	case true:
		switch streamPresence {
		case 1:
			return
		default:
			s.StreamStateMap[userID] = false
		}
	case false:
		switch streamPresence {
		case 1:
			s.StreamStateMap[userID] = true
		default:
			return
		}
	}
}

// Handler ....
func Handler(s *discordgo.Session, m *discordgo.PresenceUpdate) {
	streamHandler(s, m)
}
func streamHandler(s *discordgo.Session, p *discordgo.PresenceUpdate) {
	logger.Infof("invoking stream handler.\n session: %s\n presenceUpdate: %s", s, p)

	guildID, exists := os.LookupEnv("GUILD_ID")
	if !exists {
		logger.Error("Cannot find env variable GUILD_ID. Please ensure this is set to use gonk.")
		return
	}

	streamChannel, exists := os.LookupEnv("STREAM_CHANNEL")
	if !exists {
		logger.Error("Cannot find env variable STREAM_CHANNEL. Please ensure this is set to use streaming alerts.")
		return
	}

	if !validateGuildID(p, guildID) {
		logger.Errorf("cannot validate guild id: %s", guildID)
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
