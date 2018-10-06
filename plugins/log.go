package plugins

import (
	"github.com/evilsocket/shellz/log"
)

type logManager struct{}

func newLogManager() logManager {
	return logManager{}
}

func (l logManager) Debug(m string) {
	log.Debug("%s", m)
}

func (l logManager) Info(m string) {
	log.Info("%s", m)
}

func (l logManager) Warning(m string) {
	log.Error("%s", m)
}

func (l logManager) Error(m string) {
	log.Error("%s", m)
}
