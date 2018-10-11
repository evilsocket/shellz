package plugins

import (
	"fmt"

	"github.com/evilsocket/islazy/log"
	"github.com/evilsocket/shellz/core"
	"github.com/evilsocket/shellz/models"

	"github.com/evilsocket/islazy/async"
	"github.com/evilsocket/islazy/plugin"
)

var (
	defines = map[string]interface{}{
		"log":  getLOG(),
		"tcp":  getTCP(),
		"http": getHTTP(),
	}
)

type Plugin struct {
	*plugin.Plugin

	timeouts core.Timeouts
	ctx      interface{}
}

func LoadPlugin(path string) (error, *Plugin) {
	if p, err := plugin.Load(path, defines); err != nil {
		return err, nil
	} else {
		return nil, &Plugin{
			Plugin: p,
		}
	}
}

func (p *Plugin) NewSession(sh models.Shell, timeouts core.Timeouts) (err error, clone *Plugin) {
	p.Lock()
	defer p.Unlock()

	if err, clone = LoadPlugin(p.Path); err != nil {
		return
	}

	clone.timeouts = timeouts
	_, err = async.WithTimeout(timeouts.Connect, func() interface{} {
		clone.ctx, err = clone.Call("Create", sh)
		return err
	})
	if err != nil {
		return err, nil
	}
	return
}

func (p *Plugin) Type() string {
	return "plugin"
}

type eres struct {
	err error
	buf []byte
}

func (p *Plugin) Exec(cmd string) ([]byte, error) {
	p.Lock()
	defer p.Unlock()

	obj, err := async.WithTimeout(p.timeouts.Read+p.timeouts.Write, func() interface{} {
		if ret, err := p.Call("Exec", p.ctx, cmd); err != nil {
			return eres{err: err}
		} else if ret == nil {
			return eres{err: fmt.Errorf("return value of Exec is null")}
		} else if buf, ok := ret.([]byte); !ok {
			return eres{err: fmt.Errorf("error while converting %v to []byte", ret)}
		} else {
			return eres{buf: buf}
		}
	})
	if err != nil {
		return nil, err
	}
	er := obj.(eres)
	return er.buf, er.err
}

func (p *Plugin) Close() {
	p.Lock()
	defer p.Unlock()

	if err, _ := p.Call("Close", p.ctx); err != nil {
		log.Warning("error while running Close callback for plugin %s: %s", p.Name, err)
	}
}
