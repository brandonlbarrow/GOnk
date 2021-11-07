package stream

import (
	"os"
	"sync"
	"testing"

	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
)

var (
	log = func() *logrus.Logger {
		l := logrus.New()
		l.SetOutput(os.Stdout)
		return l
	}
)

func TestHandler_streamHandler(t *testing.T) {
	type args struct {
		s Sessioner
		p *discordgo.PresenceUpdate
	}
	tests := []struct {
		name string
		m    *Handler
		args args
		want map[string]streamerStatus
	}{
		{
			name: "success, user is streaming, change state",
			m: &Handler{
				logger:      log(),
				guildID:     "foo",
				channelID:   "bar",
				streamerMap: newStreamerMap(),
			},
			args: args{
				s: &mockSessioner{
					respGetUser: &discordgo.User{
						ID:       "foo",
						Username: "bar",
					},
					respSendMessage: &discordgo.Message{},
				},
				p: &discordgo.PresenceUpdate{
					Presence: discordgo.Presence{
						User: &discordgo.User{
							ID:       "foo",
							Username: "bar",
						},
						Game: &discordgo.Game{
							Name:    "Mystical Ninja Starring Goemon",
							URL:     "https://twitch.tv/",
							Details: "ganbare",
							Type:    discordgo.GameTypeStreaming,
						},
					},
					GuildID: "foo",
				},
			},
		},
		{
			name: "success, user is streaming but ends stream, change state",
			m: &Handler{
				logger:    log(),
				guildID:   "foo",
				channelID: "bar",
				streamerMap: &streamerMap{
					lock: sync.RWMutex{},
					streamList: map[string]streamerStatus{
						"foo": {
							statusStreamingKey: true,
						},
					},
				},
			},
			args: args{
				s: &mockSessioner{
					respGetUser: &discordgo.User{
						ID:       "foo",
						Username: "bar",
					},
					respSendMessage: &discordgo.Message{},
				},
				p: &discordgo.PresenceUpdate{
					Presence: discordgo.Presence{
						User: &discordgo.User{
							ID:       "foo",
							Username: "bar",
						},
						Game: &discordgo.Game{
							Name: "Mystical Ninja Starring Goemon",
							Type: discordgo.GameTypeGame,
						},
					},
					GuildID: "foo",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Error(tt.m.streamerMap.getStreamList())
			tt.m.streamHandler(tt.args.s, tt.args.p)
			t.Error(tt.m.streamerMap.getStreamList())
		})
	}
}

type mockSessioner struct {
	errSendMessage  error
	errGetUser      error
	respGetUser     *discordgo.User
	respSendMessage *discordgo.Message
}

func (m *mockSessioner) ChannelMessageSend(string, string) (*discordgo.Message, error) {
	return m.respSendMessage, m.errSendMessage
}

func (m *mockSessioner) User(id string) (*discordgo.User, error) {
	return m.respGetUser, m.errGetUser
}
