package plugins

import (
	"github.com/evilsocket/shellz/log"
)

type logPackage struct{}

var p = logPackage{}

func getLOG() logPackage {
	return p
}

func (l logPackage) Debug(m string) {
	log.Debug("%s", m)
}

func (l logPackage) Info(m string) {
	log.Info("%s", m)
}

func (l logPackage) Warning(m string) {
	log.Error("%s", m)
}

func (l logPackage) Error(m string) {
	log.Error("%s", m)
}
