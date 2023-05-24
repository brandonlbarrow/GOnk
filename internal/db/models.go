package db

// Server is a Discord server
type Server struct {
	// GuildID is the ID of the Discord server
	GuildID   string     `json:"guild_id"`
	Streamers []Streamer `json:"streamers"`
}

// ServerMap is a datastore mapping GuildIDs (Server IDs) to Server configurations
type ServerMap map[string]Server

// Streamer is a Discord user and channel combination to post stream updates for when receiving Twitch stream online notification events.
type Streamer struct {
	UserID    string `json:"user_id"`
	ChannelID string `json:"channel_id"`
}
