package session

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"sync"

	"github.com/robertkrimen/otto"
)

type Plugin struct {
	Name string
	Code string
	Path string

	vm       *otto.Otto
	cbCreate *otto.Script
	cbExec   *otto.Script
	cbClose  *otto.Script
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

func (p *Plugin) Clone() *Plugin {
	_, clone := LoadPlugin(p.Path, true)
	return clone
}

func (p *Plugin) Create(ctx Context) (error, otto.Value) {
	if err := p.vm.Set("ctx", ctx); err != nil {
		return err, otto.UndefinedValue()
	} else if _, err := p.vm.Run(p.cbCreate); err != nil {
		return err, otto.UndefinedValue()
	} else if ret, err := p.vm.Get("obj"); err != nil {
		return err, otto.UndefinedValue()
	} else {
		return nil, ret
	}
}

func (p *Plugin) Exec(cmd string) (error, otto.Value) {
	if err := p.vm.Set("cmd", cmd); err != nil {
		return err, otto.UndefinedValue()
	} else if _, err := p.vm.Run(p.cbExec); err != nil {
		return err, otto.UndefinedValue()
	} else if ret, err := p.vm.Get("ret"); err != nil {
		return err, otto.UndefinedValue()
	} else {
		return nil, ret
	}
}

func (p *Plugin) Close() error {
	if _, err := p.vm.Run(p.cbClose); err != nil {
		return err
	}
	return nil
}

func (p *Plugin) compile() error {
	p.vm = otto.New()
	// define built in functions and objects
	if err := p.doDefines(); err != nil {
		return err
	}
	// run the code once in order to define all the functions
	// and validate the syntax
	if _, err := p.vm.Run(p.Code); err != nil {
		return err
	}
	// validate and precompile callbacks
	if err := p.compileCall(&p.cbCreate, "Create", "var obj = Create(ctx)"); err != nil {
		return fmt.Errorf("error while compiling Create callback for %s: %s", p.Path, err)
	} else if err = p.compileCall(&p.cbExec, "Exec", "var ret = Exec(obj, cmd);"); err != nil {
		return fmt.Errorf("error while compiling Exec callback for %s: %s", p.Path, err)
	} else if err = p.compileCall(&p.cbClose, "Close", "Close(obj);"); err != nil {
		return fmt.Errorf("error while compiling Close callback for %s: %s", p.Path, err)
	}
	return nil
}

func (p *Plugin) compileCall(script **otto.Script, name string, call string) error {
	if cb, err := p.vm.Get(name); err != nil {
		return err
	} else if !cb.IsFunction() {
		return fmt.Errorf("%s is not a function", name)
	} else if *script, err = p.vm.Compile("", call); err != nil {
		return err
	}
	return nil
}

type PluginSession struct {
	sync.Mutex
	plugin *Plugin
}

func NewPluginSession(plugin *Plugin, ctx Context) (error, Session) {
	s := &PluginSession{
		plugin: plugin,
	}

	if err, _ := s.plugin.Create(ctx); err != nil {
		return err, nil
	}

	return nil, s
}

func (t *PluginSession) Type() string {
	return "plugin"
}

func (p *PluginSession) Exec(cmd string) ([]byte, error) {
	p.Lock()
	defer p.Unlock()

	if err, ret := p.plugin.Exec(cmd); err != nil {
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

func (p *PluginSession) Close() {
	p.Lock()
	defer p.Unlock()

	p.plugin.Close()
}
