package models

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/evilsocket/shellz/log"

	"github.com/evilsocket/islazy/fs"
)

var (
	Path  = "~/.shellz"
	Paths = map[string]string(nil)
)

func Init() (err error) {
	log.Debug("initializing models.Paths from %s ...", Path)

	Paths = map[string]string{
		"idents":  filepath.Join(Path, "idents"),
		"shells":  filepath.Join(Path, "shells"),
		"plugins": filepath.Join(Path, "plugins"),
	}

	for name, path := range Paths {
		if Paths[name], err = fs.Expand(path); err != nil {
			return fmt.Errorf("error while expanding path '%s': %s", path, err)
		} else {
			path = Paths[name]
		}

		log.Debug("models.Paths[%s] = %s", name, path)
		if !fs.Exists(path) {
			log.Info("creating folder %s ...", path)
			if err = os.MkdirAll(path, os.ModePerm); err != nil {
				return fmt.Errorf("error while creating path '%s': %s", path, err)
			}
		}
	}

	return
}
