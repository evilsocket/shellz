package core

const (
	Name    = "shellz"
	Version = "1.0.0"
	Author  = "Simone 'evilsocket' Margaritelli"
	Website = "https://evilsocket.net/"
)

var (
	Banner = Blue(`
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
)
