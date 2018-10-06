package plugins

func (p *Plugin) doDefines() error {
	p.vm.Set("log", getLOG())
	p.vm.Set("tcp", getTCP())
	p.vm.Set("http", getHTTP())
	return nil
}
