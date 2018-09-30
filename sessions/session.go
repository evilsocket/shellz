package sessions

import (
	"net"
)

type Handler func(address net.IP, port int, user string, pass string, keyFile string) (error, Session)

type Session interface {
	Type() string
	Exec(cmd string) ([]byte, error)
	Close()
}
