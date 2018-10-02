package core

import (
	"math/rand"
	"time"
)

const (
	Name    = "shellz"
	Version = "1.0.0"
	Author  = "Simone 'evilsocket' Margaritelli"
	Website = "https://evilsocket.net/"
)

var (
	Banner = ""
)

func init() {
	colors := []func(s string) string{
		Red,
		Blue,
		Yellow,
		Green,
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
		Dim("Made with ") + Red("❤") + Dim("  by "+Author) +
		"\n"

}
