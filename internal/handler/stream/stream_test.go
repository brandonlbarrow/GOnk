package stream

import (
	"os"
	"sync"
	"testing"

	"github.com/bwmarrin/discordgo"
	"github.com/go-test/deep"
	"github.com/sirupsen/logrus"
)

var (
	log = func() *logrus.Logger {
		l := logrus.New()
		l.SetOutput(os.Stdout)
		return l
	}
)

var streamerState = newStreamerMap()

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
				streamerMap: streamerState,
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
			want: map[string]streamerStatus{
				"foo": {
					statusStreamingKey: true,
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
			want: map[string]streamerStatus{
				"foo": {
					statusStreamingKey: false,
				},
			},
		},
		{
			name: "success, user is streaming, but got new event, still streaming",
			m: &Handler{
				logger:      log(),
				guildID:     "foo",
				channelID:   "bar",
				streamerMap: streamerState,
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
			want: map[string]streamerStatus{
				"foo": {
					statusStreamingKey: true,
				},
			},
		},
		{
			name: "success, another user is streaming, change state",
			m: &Handler{
				logger:      log(),
				guildID:     "foo",
				channelID:   "bar",
				streamerMap: streamerState,
			},
			args: args{
				s: &mockSessioner{
					respGetUser: &discordgo.User{
						ID:       "baz",
						Username: "bar",
					},
					respSendMessage: &discordgo.Message{},
				},
				p: &discordgo.PresenceUpdate{
					Presence: discordgo.Presence{
						User: &discordgo.User{
							ID:       "baz",
							Username: "bar",
						},
						Game: &discordgo.Game{
							Name:    "Super Metroid",
							URL:     "https://twitch.tv/",
							Details: "Saving the frames",
							Type:    discordgo.GameTypeStreaming,
						},
					},
					GuildID: "foo",
				},
			},
			want: map[string]streamerStatus{
				"foo": {
					statusStreamingKey: true,
				},
				"baz": {
					statusStreamingKey: true,
				},
			},
		},
		{
			name: "success, foo stopped streaming but baz continues, change state for foo",
			m: &Handler{
				logger:      log(),
				guildID:     "foo",
				channelID:   "bar",
				streamerMap: streamerState,
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
							Type:    discordgo.GameTypeGame,
						},
					},
					GuildID: "foo",
				},
			},
			want: map[string]streamerStatus{
				"foo": {
					statusStreamingKey: false,
				},
				"baz": {
					statusStreamingKey: true,
				},
			},
		},
		{
			name: "success, baz continues streaming, foo stops playing a game",
			m: &Handler{
				logger:      log(),
				guildID:     "foo",
				channelID:   "bar",
				streamerMap: streamerState,
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
					},
					GuildID: "foo",
				},
			},
			want: map[string]streamerStatus{
				"foo": {
					statusStreamingKey: false,
				},
				"baz": {
					statusStreamingKey: true,
				},
			},
		},
		{
			name: "success, baz loses power",
			m: &Handler{
				logger:      log(),
				guildID:     "foo",
				channelID:   "bar",
				streamerMap: streamerState,
			},
			args: args{
				s: &mockSessioner{
					respGetUser: &discordgo.User{
						ID:       "baz",
						Username: "bar",
					},
					respSendMessage: &discordgo.Message{},
				},
				p: &discordgo.PresenceUpdate{
					Presence: discordgo.Presence{
						User: &discordgo.User{
							ID:       "baz",
							Username: "bar",
						},
					},
					GuildID: "foo",
				},
			},
			want: map[string]streamerStatus{
				"foo": {
					statusStreamingKey: false,
				},
				"baz": {
					statusStreamingKey: false,
				},
			},
		},
		{
			name: "success single user mode, foo is the streamer",
			m: &Handler{
				logger:      log(),
				guildID:     "foo",
				channelID:   "bar",
				userID:      "foo",
				streamerMap: streamerState,
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
			want: map[string]streamerStatus{
				"foo": {
					statusStreamingKey: true,
				},
				"baz": {
					statusStreamingKey: false,
				},
			},
		},
		{
			name: "success single user mode, foo is the streamer, baz starts streaming, do not want update",
			m: &Handler{
				logger:      log(),
				guildID:     "foo",
				channelID:   "bar",
				userID:      "foo",
				streamerMap: streamerState,
			},
			args: args{
				s: &mockSessioner{
					respGetUser: &discordgo.User{
						ID:       "baz",
						Username: "bar",
					},
					respSendMessage: &discordgo.Message{},
				},
				p: &discordgo.PresenceUpdate{
					Presence: discordgo.Presence{
						User: &discordgo.User{
							ID:       "baz",
							Username: "bar",
						},
						Game: &discordgo.Game{
							Name:    "Super Metroid",
							URL:     "https://twitch.tv/",
							Details: "Saving the frames",
							Type:    discordgo.GameTypeStreaming,
						},
					},

					GuildID: "foo",
				},
			},
			want: map[string]streamerStatus{
				"foo": {
					statusStreamingKey: true,
				},
				"baz": {
					statusStreamingKey: false,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.m.streamHandler(tt.args.s, tt.args.p)
			diff := deep.Equal(tt.m.streamerMap.getStreamList(), tt.want)
			if len(diff) > 0 {
				t.Errorf("got diff %v", diff)
				return
			}
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
