package log

import (
	"os"
	"sync"
)

const (
	DEBUG = iota
	INFO
	OUTPUT
	IMPORTANT
	WARNING
	ERROR
	FATAL
)

// https://misc.flogisoft.com/bash/tip_colors_and_formatting
var (
	BOLD = "\033[1m"
	DIM  = "\033[2m"

	RED    = "\033[31m"
	GREEN  = "\033[32m"
	BLUE   = "\033[34m"
	YELLOW = "\033[33m"

	FG_BLACK = "\033[30m"
	FG_WHITE = "\033[97m"

	BG_DGRAY  = "\033[100m"
	BG_RED    = "\033[41m"
	BG_GREEN  = "\033[42m"
	BG_YELLOW = "\033[43m"
	BG_LBLUE  = "\033[104m"

	RESET = "\033[0m"
)

var (
	NoColors      = false
	DebugMessages = false
	File          = ""

	labels = map[int]string{
		DEBUG:     "dbg",
		INFO:      "inf",
		OUTPUT:    "out",
		IMPORTANT: "imp",
		WARNING:   "war",
		ERROR:     "err",
		FATAL:     "!!!",
	}

	colors = map[int]string{
		DEBUG:     DIM + FG_BLACK + BG_DGRAY,
		INFO:      FG_WHITE + BG_GREEN,
		OUTPUT:    DIM + FG_BLACK + BG_DGRAY,
		IMPORTANT: FG_WHITE + BG_LBLUE,
		WARNING:   FG_WHITE + BG_YELLOW,
		ERROR:     FG_WHITE + BG_RED,
		FATAL:     FG_WHITE + BG_RED + BOLD,
	}

	lock   = &sync.Mutex{}
	writer = os.Stdout
)

func Init() {
	if File != "" {
		var err error
		if writer, err = os.OpenFile(File, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644); err != nil {
			panic(err)
		}
	}

	if NoColors {
		for level := range colors {
			colors[level] = ""
		}

		BOLD = ""
		DIM = ""
		RED = ""
		GREEN = ""
		BLUE = ""
		YELLOW = ""
		FG_BLACK = ""
		FG_WHITE = ""
		BG_DGRAY = ""
		BG_RED = ""
		BG_GREEN = ""
		BG_YELLOW = ""
		BG_LBLUE = ""
		RESET = ""
	}
}

func Close() {
	if writer != os.Stdout {
		writer.Close()
	}
}
