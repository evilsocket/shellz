package session

import (
	"fmt"
)

type Proxy struct {
	Address  string
	Port     int
	Username string
	Password string
}

func (p Proxy) String() string {
	if p.Username != "" {
		return fmt.Sprintf("%s:%s@%s:%d", p.Username, p.Password, p.Address, p.Port)
	}
	return fmt.Sprintf("%s:%d", p.Address, p.Port)
}
