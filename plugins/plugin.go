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

	vm       *otto.Otto
	cbCreate otto.Value
	cbExec   otto.Value
	cbClose  otto.Value
	ctx      otto.Value
}

func LoadPlugin(path string, doCompile bool) (error, *Plugin) {
	raw, err := ioutil.ReadFile(path)
	if err != nil {
		return err, nil
	}

	plugin := &Plugin{
		Name: strings.Replace(filepath.Base(path), ".js", "", -1),
		Code: string(raw),
		Path: path,
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
	} else if p.cbCreate, err = p.findCall("Create"); err != nil {
		return fmt.Errorf("error while compiling Create callback for %s: %s", p.Path, err)
	} else if p.cbExec, err = p.findCall("Exec"); err != nil {
		return fmt.Errorf("error while compiling Exec callback for %s: %s", p.Path, err)
	} else if p.cbClose, err = p.findCall("Close"); err != nil {
		return fmt.Errorf("error while compiling Close callback for %s: %s", p.Path, err)
	}
	return nil
}

func (p *Plugin) clone() *Plugin {
	_, clone := LoadPlugin(p.Path, true)
	return clone
}

func (p *Plugin) NewSession(sh models.Shell, timeouts core.Timeouts) (err error, clone *Plugin) {
	p.Lock()
	defer p.Unlock()
	clone = p.clone()
	clone.ctx, err = clone.cbCreate.Call(otto.NullValue(), sh)
	return
}

func (p *Plugin) Type() string {
	return "plugin"
}

func (p *Plugin) Exec(cmd string) ([]byte, error) {
	p.Lock()
	defer p.Unlock()

	if ret, err := p.cbExec.Call(otto.NullValue(), p.ctx, cmd); err != nil {
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

func (p *Plugin) Close() {
	p.Lock()
	defer p.Unlock()
	if _, err := p.cbClose.Call(otto.NullValue(), p.ctx); err != nil {
		log.Warning("error while running Close callback for plugin %s: %s", p.Name, err)
	}
}
