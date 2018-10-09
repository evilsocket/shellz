package models

import (
	"fmt"
	"net"
	"strconv"
)

type Proxy struct {
	Address  string `json:"address"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func (p Proxy) Empty() bool {
	return p.Address == ""
}

func (p Proxy) String() string {
	host := net.JoinHostPort(p.Address, strconv.Itoa(p.Port))
	if p.Username != "" {
		return fmt.Sprintf("%s:%s@%s", p.Username, p.Password, host)
	}
	return host
}
