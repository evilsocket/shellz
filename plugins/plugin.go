package plugins

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"sync"

	"github.com/evilsocket/shellz/core"
	"github.com/evilsocket/shellz/log"
	"github.com/evilsocket/shellz/models"

	"github.com/robertkrimen/otto"
)

type Plugin struct {
	sync.Mutex

	Name string
	Code string
	Path string

	vm        *otto.Otto
	callbacks map[string]otto.Value
	ctx       interface{}
}

func LoadPlugin(path string, doCompile bool) (error, *Plugin) {
	raw, err := ioutil.ReadFile(path)
	if err != nil {
		return err, nil
	}

	plugin := &Plugin{
		Name:      strings.Replace(filepath.Base(path), ".js", "", -1),
		Code:      string(raw),
		Path:      path,
		callbacks: make(map[string]otto.Value),
	}

	if doCompile {
		if err = plugin.compile(); err != nil {
			return err, nil
		}
	}
	return nil, plugin
}

func (p *Plugin) findCall(name string) (cb otto.Value, err error) {
	if cb, err = p.vm.Get(name); err != nil {
		return
	} else if !cb.IsFunction() {
		err = fmt.Errorf("%s is not a function", name)
	}
	return
}

func (p *Plugin) compile() (err error) {
	p.vm = otto.New()
	// define built in functions and objects
	if err = p.doDefines(); err != nil {
		return
	}
	// run the code once in order to define all the functions
	// and validate the syntax, then get the callbacks
	if _, err = p.vm.Run(p.Code); err != nil {
		return
	} else if p.callbacks["Create"], err = p.findCall("Create"); err != nil {
		return fmt.Errorf("error while compiling Create callback for %s: %s", p.Path, err)
	} else if p.callbacks["Exec"], err = p.findCall("Exec"); err != nil {
		return fmt.Errorf("error while compiling Exec callback for %s: %s", p.Path, err)
	} else if p.callbacks["Close"], err = p.findCall("Close"); err != nil {
		return fmt.Errorf("error while compiling Close callback for %s: %s", p.Path, err)
	}
	return nil
}

func (p *Plugin) clone() *Plugin {
	_, clone := LoadPlugin(p.Path, true)
	return clone
}

func (p *Plugin) call(name string, args ...interface{}) (error, interface{}) {
	if cb, found := p.callbacks[name]; !found {
		return fmt.Errorf("%s does not name a function", name), nil
	} else if ret, err := cb.Call(otto.NullValue(), args...); err != nil {
		return err, nil
	} else if !ret.IsUndefined() {
		if exported, err := ret.Export(); err != nil {
			return err, nil
		} else {
			return nil, exported
		}
	}
	return nil, nil
}

func (p *Plugin) NewSession(sh models.Shell, timeouts core.Timeouts) (err error, clone *Plugin) {
	p.Lock()
	defer p.Unlock()
	clone = p.clone()
	err, clone.ctx = clone.call("Create", sh)
	return
}

func (p *Plugin) Type() string {
	return "plugin"
}

func (p *Plugin) Exec(cmd string) ([]byte, error) {
	p.Lock()
	defer p.Unlock()

	if err, ret := p.call("Exec", p.ctx, cmd); err != nil {
		return nil, err
	} else if ret == nil {
		return nil, fmt.Errorf("return value of Exec is null")
	} else if array, ok := ret.([]byte); !ok {
		return nil, fmt.Errorf("error while converting %v to []byte", ret)
	} else {
		return array, nil
	}
}

func (p *Plugin) Close() {
	p.Lock()
	defer p.Unlock()
	if err, _ := p.call("Close", p.ctx); err != nil {
		log.Warning("error while running Close callback for plugin %s: %s", p.Name, err)
	}
}
