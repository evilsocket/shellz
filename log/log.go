package log

import (
	"fmt"
	"os"
	"sync"

	"github.com/evilsocket/shellz/core"
)

const (
	DEBUG = iota
	INFO
	IMPORTANT
	WARNING
	ERROR
	FATAL
)

var (
	DebugMessages = false
	lock          = &sync.Mutex{}

	Labels = map[int]string{
		DEBUG:     "dbg",
		INFO:      "inf",
		IMPORTANT: "imp",
		WARNING:   "war",
		ERROR:     "err",
		FATAL:     "!!!",
	}

	Colors = map[int]string{
		DEBUG:     core.DIM + core.FG_BLACK + core.BG_DGRAY,
		INFO:      core.FG_WHITE + core.BG_GREEN,
		IMPORTANT: core.FG_WHITE + core.BG_LBLUE,
		WARNING:   core.FG_WHITE + core.BG_YELLOW,
		ERROR:     core.FG_WHITE + core.BG_RED,
		FATAL:     core.FG_WHITE + core.BG_RED + core.BOLD,
	}
)

func Raw(format string, args ...interface{}) {
	lock.Lock()
	defer lock.Unlock()

	fmt.Fprintf(os.Stdout, format, args...)
	fmt.Fprintf(os.Stdout, "\n")
}

func color(level int, format string, args ...interface{}) {
	label := Labels[level]
	color := Colors[level]
	Raw(color+label+core.RESET+" "+format, args...)
}

func Info(format string, args ...interface{}) {
	color(INFO, format, args...)
}

func Warning(format string, args ...interface{}) {
	color(WARNING, format, args...)
}

func Error(format string, args ...interface{}) {
	color(ERROR, format, args...)
}

func Fatal(format string, args ...interface{}) {
	color(FATAL, format, args...)
	os.Exit(1)
}

func Debug(format string, args ...interface{}) {
	if DebugMessages {
		color(DEBUG, format, args...)
	}
}
