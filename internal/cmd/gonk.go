package cmd

import (
	"fmt"
	"strings"
)

type RootCmd struct {
	cmd string
}

func NewRootCmd(command string) RootCmd {
	return RootCmd{cmd: command}
}

func (r RootCmd) Prefix(command string) string {
	if strings.HasPrefix(command, r.cmd) {
		return command
	} else {
		return fmt.Sprintf("%s %s", r.cmd, command)
	}
}
