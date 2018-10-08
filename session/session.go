package session

import (
	"time"
)

type Timeouts struct {
	Connect time.Duration
	Read    time.Duration
	Write   time.Duration
}

type Context struct {
	Host     string
	Port     int
	Username string
	Password string
	KeyFile  string
	Ciphers  []string
	Timeouts Timeouts
	Proxy    Proxy
}

type Handler func(ctx Context) (error, Session)

type Session interface {
	Type() string
	Exec(cmd string) ([]byte, error)
	Close()
}
