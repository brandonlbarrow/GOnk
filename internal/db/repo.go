package db

// ServerRepo defines the data interface for Gonk's Server data
type ServerRepo interface {
	GetServerByID(guildID string) (*Server, error)
	GetServers() (map[string]*Server, error)
	AddServer(guildID string, server *Server) error
	UpdateServer(guildID string, server *Server) error
	DeleteServer(guildID string) error
}
