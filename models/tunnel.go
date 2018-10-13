package models

import (
	"fmt"
	"net"
	"strconv"

	"github.com/evilsocket/islazy/tui"
)

type Address struct {
	Address string `json:"address"`
	Port    int    `json:"port"`
}

func (t Address) String() string {
	return net.JoinHostPort(t.Address, strconv.Itoa(t.Port))
}

func (t Address) Empty() bool {
	return t.Address == ""
}

type Tunnel struct {
	Local  Address `json:"local"`
	Remote Address `json:"remote"`
}

func (t Tunnel) Empty() bool {
	return t.Local.Empty() || t.Remote.Empty()
}

func (t Tunnel) String() string {
	if t.Empty() {
		return ""
	}
	return fmt.Sprintf("%s %s %s", t.Local.String(), tui.Dim("<->"), t.Remote.String())
}
