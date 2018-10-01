package models

import (
	"os"

	"github.com/evilsocket/shellz/core"
	"github.com/evilsocket/shellz/log"
)

var (
	Paths = map[string]string{
		"base":    "~/.shellz",
		"idents":  "~/.shellz/idents",
		"shells":  "~/.shellz/shells",
		"plugins": "~/.shellz/plugins",
	}
)

func init() {
	var err error
	for name, path := range Paths {
		if Paths[name], err = core.ExpandPath(path); err != nil {
			log.Fatal("error while expanding path '%s': %s", path, err)
		} else {
			path = Paths[name]
		}

		log.Debug("models.Paths[%s] = %s", name, path)
		if !core.Exists(path) {
			log.Info("creating folder %s ...", path)
			if err = os.MkdirAll(path, os.ModePerm); err != nil {
				log.Fatal("error while creating path '%s': %s", path, err)
			}
		}
	}
}
