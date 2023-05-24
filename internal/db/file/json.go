package file

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sync"

	"github.com/brandonlbarrow/gonk/v2/internal/db"
)

type JSONFileDB struct {
	db db.ServerMap
	m  *sync.Mutex
}

func (j *JSONFileDB) GetServerByID(guildID string) (*db.Server, error) {
	panic("not implemented") // TODO: Implement
}

func (j *JSONFileDB) GetServers() (map[string]*db.Server, error) {
	panic("not implemented") // TODO: Implement
}

func (j *JSONFileDB) AddServer(guildID string, server *db.Server) error {
	panic("not implemented") // TODO: Implement
}

func (j *JSONFileDB) UpdateServer(guildID string, server *db.Server) error {
	panic("not implemented") // TODO: Implement
}

func (j *JSONFileDB) DeleteServer(guildID string) error {
	panic("not implemented") // TODO: Implement
}

func NewJSONFileDB(path string) (*JSONFileDB, error) {
	s, err := os.Stat(path)
	if err != nil {
		workingDir, getWdErr := os.Getwd()
		if getWdErr != nil {
			return nil, err
		} else {
			return nil, fmt.Errorf("could not find provided path in current directory: %s. Error: %w", workingDir, err)

		}
	}
	if !s.IsDir() {
		contents, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}
		var serverMap db.ServerMap
		if err := json.Unmarshal(contents, &serverMap); err != nil {
			return nil, err
		}
		m := &sync.Mutex{}
		return &JSONFileDB{
			m:  m,
			db: serverMap,
		}, nil

	} else {
		return nil, errors.New("not a full path to a file")
	}
}
