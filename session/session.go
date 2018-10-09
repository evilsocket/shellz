package session

import (
	"github.com/evilsocket/shellz/core"
	"github.com/evilsocket/shellz/models"
)

var (
	managers = map[string]Handler{
		"ssh":    NewSSH,
		"telnet": NewTelnet,
	}
)

type Handler func(sh models.Shell, timeouts core.Timeouts) (error, Session)

type Session interface {
	Type() string
	Exec(cmd string) ([]byte, error)
	Close()
}

func For(sh models.Shell, timeouts core.Timeouts) (error, Session) {
	if mgr := managers[sh.Type]; mgr != nil {
		return mgr(sh, timeouts)
	}
	return nil, nil
}
