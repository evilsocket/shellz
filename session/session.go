package session

import (
	"net"
)

type Context struct {
	Address  net.IP
	Port     int
	Username string
	Password string
	KeyFile  string
}

type Handler func(ctx Context) (error, Session)

type Session interface {
	Type() string
	Exec(cmd string) ([]byte, error)
	Close()
}
