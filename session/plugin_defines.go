package session

import (
	"github.com/evilsocket/shellz/log"

	"github.com/robertkrimen/otto"
)

func (p *Plugin) doDefines() error {
	p.vm.Set("http", newHttpClient())
	p.vm.Set("log", func(call otto.FunctionCall) otto.Value {
		for _, v := range call.ArgumentList {
			log.Info("%s", v.String())
		}
		return otto.Value{}
	})

	return nil
}
