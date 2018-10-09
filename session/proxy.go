package session

import (
	"fmt"
	"net"
	"strconv"
)

type Proxy struct {
	Address  string
	Port     int
	Username string
	Password string
}

func (p Proxy) String() string {
	host := net.JoinHostPort(p.Address, strconv.Itoa(p.Port))
	if p.Username != "" {
		return fmt.Sprintf("%s:%s@%s", p.Username, p.Password, host)
	}
	return host
}
