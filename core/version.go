package core

import (
	"math/rand"
	"time"

	"github.com/evilsocket/islazy/tui"
)

const (
	Name    = "shellz"
	Version = "1.5.1"
	Author  = "Simone 'evilsocket' Margaritelli"
	Website = "https://evilsocket.net/"
)

var (
	Banner = ""
)

func init() {
	colors := []func(s string) string{
		tui.Red,
		tui.Blue,
		tui.Yellow,
		tui.Green,
	}
	rand.Seed(time.Now().Unix())

	Banner = colors[rand.Intn(len(colors))](`
  ██████  ██░ ██ ▓█████  ██▓     ██▓    ▒███████▒
▒██    ▒ ▓██░ ██▒▓█   ▀ ▓██▒    ▓██▒    ▒ ▒ ▒ ▄▀░
░ ▓██▄   ▒██▀▀██░▒███   ▒██░    ▒██░    ░ ▒ ▄▀▒░ 
  ▒   ██▒░▓█ ░██ ▒▓█  ▄ ▒██░    ▒██░      ▄▀▒   ░
▒██████▒▒░▓█▒░██▓░▒████▒░██████▒░██████▒▒███████▒
▒ ▒▓▒ ▒ ░ ▒ ░░▒░▒░░ ▒░ ░░ ▒░▓  ░░ ▒░▓  ░░▒▒ ▓░▒░▒
░ ░▒  ░ ░ ▒ ░▒░ ░ ░ ░  ░░ ░ ▒  ░░ ░ ▒  ░░░▒ ▒ ░ ▒
░  ░  ░   ░  ░░ ░   ░     ░ ░     ░ ░   ░ ░ ░ ░ ░
      ░   ░  ░  ░   ░  ░    ░  ░    ░  ░  ░ ░    
                                        ░  `) +
		"v" + Version + "\n" +
		tui.Dim("Made with ") + tui.Red("❤") + tui.Dim("  by "+Author) +
		"\n"

}
