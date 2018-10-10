package plugins

import (
	"fmt"

	"github.com/evilsocket/islazy/log"
	"github.com/evilsocket/shellz/models"

	"github.com/evilsocket/islazy/fs"
)

type Plugins map[string]*Plugin

var (
	plugins = make(Plugins)
)

func Load(path string) error {
	log.Debug("loading plugins from %s ...", path)
	return fs.Glob(path, "*.js", func(fileName string) error {
		if err, plugin := LoadPlugin(fileName); err != nil {
			return fmt.Errorf("error while loading plugin '%s': %s", fileName, err)
		} else if taken, found := plugins[plugin.Name]; found {
			return fmt.Errorf("plugin '%s' has name %s which is already taken by '%s'", fileName, plugin.Name, taken.Path)
		} else {
			plugins[plugin.Name] = plugin
		}
		return nil
	})
}

func Get(sh models.Shell) *Plugin {
	return plugins[sh.Type]
}

func Number() int {
	return len(plugins)
}

func Each(cb func(p *Plugin)) {
	for _, p := range plugins {
		cb(p)
	}
}
