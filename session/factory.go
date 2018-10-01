package session

import (
	"fmt"
	"path/filepath"

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

	if files, err := filepath.Glob(filepath.Join(path, "*.js")); err != nil {
		return err
	} else {
		for _, fileName := range files {
			if err, plugin := LoadPlugin(fileName, true); err != nil {
				return fmt.Errorf("error while loading plugin '%s': %s", fileName, err)
			} else if taken, found := plugins[plugin.Name]; found {
				return fmt.Errorf("plugin '%s' has name %s which is already taken by '%s'", fileName, plugin.Name, taken.Path)
			} else {
				plugins[plugin.Name] = plugin
			}
		}
	}

	log.Debug("plugins = %v", plugins)

	return nil
}

func GetManager(name string) Handler {
	return managers[name]
}

func GetPlugin(name string) *Plugin {
	return plugins[name]
}
