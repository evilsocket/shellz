package log

import (
	"fmt"
	"os"
)

func Raw(format string, args ...interface{}) {
	lock.Lock()
	defer lock.Unlock()

	fmt.Fprintf(writer, format, args...)
	fmt.Fprintf(writer, "\n")
}

func color(level int, format string, args ...interface{}) {
	Raw(colors[level]+labels[level]+RESET+" "+format, args...)
}

func Info(format string, args ...interface{}) {
	color(INFO, format, args...)
}

func Output(format string, args ...interface{}) {
	color(OUTPUT, format, args...)
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
