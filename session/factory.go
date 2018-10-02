package session

import (
	"fmt"

	"github.com/evilsocket/shellz/core"
	"github.com/evilsocket/shellz/log"
)

type Plugins map[string]*Plugin

var (
	managers = map[string]Handler{
		"ssh":    NewSSH,
		"telnet": NewTelnet,
	}

	plugins = make(Plugins)
)

func LoadPlugins(path string) error {
	log.Info("loading plugins from %s ...", path)
	return core.Glob(path, "*.js", func(fileName string) error {
		if err, plugin := LoadPlugin(fileName, true); err != nil {
			return fmt.Errorf("error while loading plugin '%s': %s", fileName, err)
		} else if taken, found := plugins[plugin.Name]; found {
			return fmt.Errorf("plugin '%s' has name %s which is already taken by '%s'", fileName, plugin.Name, taken.Path)
		} else {
			plugins[plugin.Name] = plugin
		}
		return nil
	})
}

func GetManager(name string) Handler {
	return managers[name]
}

func GetPlugin(name string) *Plugin {
	return plugins[name]
}
