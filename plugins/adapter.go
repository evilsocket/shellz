package plugins

import (
	"fmt"
	"sync"

	"github.com/evilsocket/shellz/session"
)

type SessionAdapter struct {
	sync.Mutex
	plugin *Plugin
}

func ForSession(plugin *Plugin, ctx session.Context) (error, *SessionAdapter) {
	adapter := &SessionAdapter{plugin: plugin}
	if err, _ := adapter.plugin.Create(ctx); err != nil {
		return err, nil
	}
	return nil, adapter
}

func (a *SessionAdapter) Type() string {
	return "plugin"
}

func (a *SessionAdapter) Exec(cmd string) ([]byte, error) {
	a.Lock()
	defer a.Unlock()

	if err, ret := a.plugin.Exec(cmd); err != nil {
		return nil, err
	} else if ret.IsNull() || ret.IsUndefined() {
		return []byte{}, nil
	} else if exported, err := ret.Export(); err != nil {
		return nil, err
	} else if array, ok := exported.([]byte); !ok {
		return nil, fmt.Errorf("error while converting %v to []byte", exported)
	} else {
		return array, nil
	}
}

func (a *SessionAdapter) Close() {
	a.Lock()
	defer a.Unlock()
	a.plugin.Close()
}
