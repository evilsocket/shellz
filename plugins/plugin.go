package plugins

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"sync"
	"unicode"

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

	timeouts  core.Timeouts
	vm        *otto.Otto
	callbacks map[string]otto.Value
	objects   map[string]otto.Value
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
		objects:   make(map[string]otto.Value),
	}

	if doCompile {
		if err = plugin.compile(); err != nil {
			return err, nil
		}
	}
	return nil, plugin
}

func (p *Plugin) compile() (err error) {
	p.vm = otto.New()

	// define built in functions and objects
	if err = p.doDefines(); err != nil {
		return
	}

	// track objects already defined by Otto
	predefined := map[string]bool{}
	for name := range p.vm.Context().Symbols {
		predefined[name] = true
	}

	// run the code once in order to define all the functions
	// and validate the syntax, then get the callbacks
	if _, err = p.vm.Run(p.Code); err != nil {
		return
	}

	// every uppercase object is considered exported
	for name, sym := range p.vm.Context().Symbols {
		// ignore predefined objects
		if _, found := predefined[name]; !found {
			// ignore lowercase global objects
			if unicode.IsUpper(rune(name[0])) {
				if sym.IsFunction() {
					p.callbacks[name] = sym
				} else {
					p.objects[name] = sym
				}
			}
		}
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
	clone.timeouts = timeouts
	err, obj := core.WithTimeout(timeouts.Connect, func() interface{} {
		err, clone.ctx = clone.call("Create", sh)
		return err
	})
	if err != nil {
		return err, nil
	} else if obj != nil {
		if err = obj.(error); err != nil {
			return err, nil
		}
	}
	return
}

func (p *Plugin) Type() string {
	return "plugin"
}

type eres struct {
	err   error
	array []byte
}

func (p *Plugin) Exec(cmd string) ([]byte, error) {
	p.Lock()
	defer p.Unlock()

	err, obj := core.WithTimeout(p.timeouts.Read+p.timeouts.Write, func() interface{} {
		if err, ret := p.call("Exec", p.ctx, cmd); err != nil {
			return eres{err: err}
		} else if ret == nil {
			return eres{err: fmt.Errorf("return value of Exec is null")}
		} else if array, ok := ret.([]byte); !ok {
			return eres{err: fmt.Errorf("error while converting %v to []byte", ret)}
		} else {
			return eres{array: array}
		}
	})
	if err != nil {
		return nil, err
	}
	er := obj.(eres)
	return er.array, er.err
}

func (p *Plugin) Close() {
	p.Lock()
	defer p.Unlock()
	if err, _ := p.call("Close", p.ctx); err != nil {
		log.Warning("error while running Close callback for plugin %s: %s", p.Name, err)
	}
}
