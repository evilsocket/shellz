package log

import (
	"os"
	"sync"

	"github.com/evilsocket/shellz/core"
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

var (
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
		DEBUG:     core.DIM + core.FG_BLACK + core.BG_DGRAY,
		INFO:      core.FG_WHITE + core.BG_GREEN,
		OUTPUT:    core.DIM + core.FG_BLACK + core.BG_DGRAY,
		IMPORTANT: core.FG_WHITE + core.BG_LBLUE,
		WARNING:   core.FG_WHITE + core.BG_YELLOW,
		ERROR:     core.FG_WHITE + core.BG_RED,
		FATAL:     core.FG_WHITE + core.BG_RED + core.BOLD,
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
}

func Close() {
	if wrter != os.Stdout {
		writer.Close()
	}
}
