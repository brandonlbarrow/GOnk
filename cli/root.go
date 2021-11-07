package cli

import (
	"github.com/brandonlbarrow/gonk/internal/listener"
)

type Command interface {
	exec() error
}

type RootCommand struct {
	g *listener.Listener
}
